package art.ameliah.libs.weave.leaderboard;

import art.ameliah.libs.weave.Result;
import art.ameliah.libs.weave.Utils;
import art.ameliah.libs.weave.WeaveException;
import com.google.gson.JsonArray;
import com.google.gson.JsonElement;
import com.google.gson.JsonObject;
import org.apache.http.impl.client.CloseableHttpClient;

import javax.annotation.Nullable;
import java.util.HashMap;

/**
 * Leaderboard factory
 */
public class LeaderboardFactory implements ILeaderboardFactory {

    private final HashMap<String, Leaderboard> converter = new HashMap<>();

    /**
     * Internal constructor
     */
    protected LeaderboardFactory(String baseURL, CloseableHttpClient httpClient) throws WeaveException {
        String url = String.format("%s/leaderboard_api/games/true", baseURL);
        Result<JsonArray, WeaveException> result = Utils.tryContentStringWithJsonEncoding(url, httpClient);
        if (result.isErr()) {
            throw WeaveException.fromResult(result);
        }

        for (JsonElement el : result.getValue()) {
            JsonObject obj = el.getAsJsonObject();
            Leaderboard game = new Leaderboard(obj.get("game").getAsString(), obj.get("display_name").getAsString(), true, obj.get("score_type").getAsString());
            for (JsonElement aliasElement : obj.get("aliases").getAsJsonArray()) {
                converter.put(aliasElement.getAsString(), game);
            }
            converter.put(game.name(), game);
            converter.put(game.displayName(), game);
        }
    }

    /**
     * Internal method
     */
    @Override
    public @Nullable Leaderboard getLeaderboard(String name) {
        return converter.get(name);
    }

}
