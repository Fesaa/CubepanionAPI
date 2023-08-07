package art.ameliah.libs.weave;

/**
 * Weave specific exception, any Exception emitted by Weave after construction will be of this type
 */
public class WeaveException extends Exception {

    /**
     * Artificial exception
     *
     * @param msg extra info
     */
    protected WeaveException(String msg) {
        super(msg);
    }

    /**
     * Artificial wrapper exception
     *
     * @param msg   extra info
     * @param cause original exception
     */
    protected WeaveException(String msg, Throwable cause) {
        super(msg, cause);
    }

    /**
     * Wrapper exception
     *
     * @param cause original exception
     */
    protected WeaveException(Throwable cause) {
        super(cause);
    }

}
