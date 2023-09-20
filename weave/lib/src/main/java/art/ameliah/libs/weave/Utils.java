package art.ameliah.libs.weave;

import com.google.gson.Gson;
import org.asynchttpclient.AsyncHttpClient;

import java.util.concurrent.CompletableFuture;

class Utils {

    static <T> CompletableFuture<T> makeRequest(AsyncHttpClient httpClient, String url, Class<T> clazz) {
        return httpClient
                .prepareGet(url)
                .execute()
                .toCompletableFuture()
                .thenApplyAsync(response -> {
                    WeaveException e = WeaveException.fromResponse(response, 200, true);
                    if (e != null) {
                        throw new RuntimeException(e);
                    }

                    String json = response.getResponseBody();
                    return (new Gson()).fromJson(json, clazz);
                });
    }
}
