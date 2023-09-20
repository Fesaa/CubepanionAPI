package art.ameliah.libs.weave;

import org.asynchttpclient.Response;

import javax.annotation.Nullable;


/**
 * Weave specific exception, any Exception emitted by Weave after construction will be of this type
 */
public class WeaveException extends Exception {

    /**
     * Artificial exception
     *
     * @param msg extra info
     */
    WeaveException(String msg) {
        super(msg);
    }

    /**
     * Artificial wrapper exception
     *
     * @param msg   extra info
     * @param cause original exception
     */
    WeaveException(String msg, Throwable cause) {
        super(msg, cause);
    }

    /**
     * Wrapper exception
     *
     * @param cause original exception
     */
    WeaveException(Throwable cause) {
        super(cause);
    }

    static @Nullable WeaveException fromResponse(Response response, int wantedStatus, boolean needsBody) {
        if (response == null) {
            return new WeaveException("Failed to get a response");
        }

        if (response.hasResponseHeaders()) {
            String rate = response.getHeader("retry-after");
            if (rate != null) {
                return new RateLimited(rate);
            }
        }

        if (response.getStatusCode() != wantedStatus) {
            String msg = String.format("Got %d as status code, wanted %d.\n%s",
                    wantedStatus, response.getStatusCode(),
                    response.hasResponseBody() ? response.getResponseBody() : "Unknown.");
            return new WeaveException(msg);
        }

        if (needsBody && !response.hasResponseBody()) {
            return new WeaveException("Wanted a response body, but did not get one");
        }

        return null;
    }

}
