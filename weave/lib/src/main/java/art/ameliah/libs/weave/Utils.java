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

class Utils {

    static JsonArray tryContentStringWithJsonEncoding(String url, CloseableHttpClient httpClient) throws WeaveException {
        HttpGet req = new HttpGet(url);
        HttpResponse response;
        try {
            response = httpClient.execute(req);
        } catch (IOException e) {
            throw new WeaveException("Could not execute request", e);
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
            throw new WeaveException("Rate limited, try again in: " + again + "s");
        }

        HttpEntity entity = response.getEntity();
        if (entity == null) {
            throw new WeaveException("Response entity was null");
        }

        Header contentTypeHeader = entity.getContentType();
        if (contentTypeHeader == null) {
            throw new WeaveException("No content type");
        }

        String contentType = contentTypeHeader.getValue();
        if (!contentType.contains("application/json")) {


            throw new WeaveException("Responded with non json. API error");
        }

        String responseBody;
        try {
            responseBody = EntityUtils.toString(entity);
        } catch (IOException e) {
            throw new WeaveException("Could not get responseBody", e);
        }
        try {
            return new Gson().fromJson(responseBody, JsonArray.class);
        } catch (JsonSyntaxException e) {
            throw new WeaveException("Could not construct JsonArray", e);
        }
    }

}
