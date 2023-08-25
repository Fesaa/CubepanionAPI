package art.ameliah.libs.weave;

import com.google.gson.Gson;
import com.google.gson.JsonArray;
import com.google.gson.JsonSyntaxException;
import org.apache.http.Header;
import org.apache.http.HttpEntity;
import org.apache.http.HttpResponse;
import org.apache.http.client.methods.HttpGet;
import org.apache.http.impl.client.CloseableHttpClient;
import org.apache.http.util.EntityUtils;

import java.io.IOException;

public class Utils {

    public static Result<JsonArray, WeaveException> tryContentStringWithJsonEncoding(String url, CloseableHttpClient httpClient) {
        return tryContentStringWithJsonEncoding(url, httpClient, JsonArray.class);
    }

    public static <T> Result<T, WeaveException> tryContentStringWithJsonEncoding(String url, CloseableHttpClient httpClient, Class<T> clazz) {
        HttpGet req = new HttpGet(url);
        HttpResponse response;
        try {
            response = httpClient.execute(req);
        } catch (IOException e) {
            return Result.Err(new WeaveException("Could not execute request", e));
        }

        if (response.containsHeader("retry-after")) {
            Header[] headers = response.getHeaders("retry-after");
            String again;
            if (headers.length > 0) {
                Header header = headers[0];
                again = header.getValue();
            } else {
                again = "60";
            }
            return Result.Err(new WeaveException("Rate limited, try again in: " + again + "s"));
        }

        HttpEntity entity = response.getEntity();
        if (entity == null) {
            return Result.Err(new WeaveException("Response entity was null"));
        }

        Header contentTypeHeader = entity.getContentType();
        if (contentTypeHeader == null) {
            return Result.Err(new WeaveException("No content type"));
        }

        String contentType = contentTypeHeader.getValue();
        if (!contentType.contains("application/json")) {
            return Result.Err(new WeaveException("Responded with non json. API error"));
        }

        String responseBody;
        try {
            responseBody = EntityUtils.toString(entity);
        } catch (IOException e) {
            return Result.Err(new WeaveException("Could not get responseBody", e));
        }
        try {
            return Result.Ok(new Gson().fromJson(responseBody, clazz));
        } catch (JsonSyntaxException e) {
            return Result.Err(new WeaveException("Could not construct JsonArray", e));
        }
    }

}
