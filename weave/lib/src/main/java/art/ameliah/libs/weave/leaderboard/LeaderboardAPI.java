package art.ameliah.libs.weave.leaderboard;

import art.ameliah.libs.weave.Result;
import art.ameliah.libs.weave.Utils;
import art.ameliah.libs.weave.WeaveException;
import com.google.gson.Gson;
import com.google.gson.JsonArray;
import com.google.gson.JsonElement;
import com.google.gson.JsonObject;
import org.apache.http.HttpEntity;
import org.apache.http.HttpResponse;
import org.apache.http.StatusLine;
import org.apache.http.client.methods.HttpPost;
import org.apache.http.entity.ContentType;
import org.apache.http.entity.StringEntity;
import org.apache.http.impl.client.CloseableHttpClient;
import org.apache.http.util.EntityUtils;

import java.io.IOException;
import java.util.ArrayList;
import java.util.List;
import java.util.Set;
import java.util.UUID;

/**
 * Interaction class for LeaderboardAPI
 */
public class LeaderboardAPI {

    private final String baseURL;

    private final CloseableHttpClient httpClient;

    private final ILeaderboardFactory factory;

    /**
     * @param url        base API url
     * @param httpClient connection client
     */
    public LeaderboardAPI(String url, CloseableHttpClient httpClient) {
        this.baseURL = url;
        this.httpClient = httpClient;

        ILeaderboardFactory factory1;
        try {
           factory1 = new LeaderboardFactory(url, httpClient);
        } catch (WeaveException e) {
            factory1 = new LeaderboardFactoryDummy();
        }
        this.factory = factory1;
    }

    /**
     * Submit a new leaderboard
     *
     * @param uuid    Submitters uuid
     * @param game    name
     * @param entries LeaderboardRows
     * @return StatusCode or Error
     */
    public Result<Integer, WeaveException> submitLeaderboard(UUID uuid, Leaderboard game, Set<LeaderboardRow> entries) {
        if (entries.size() != 200) {
            return Result.Err(new WeaveException("entries set must contain exactly 200 entries."));
        }
        if (!game.active()) {
            return Result.Err(new WeaveException("Game is not active, cannot submit"));
        }

        JsonObject main = new JsonObject();
        main.addProperty("uuid", uuid.toString());
        main.addProperty("unix_time_stamp", System.currentTimeMillis());
        main.addProperty("game", game.displayName());
        JsonArray cachedEntries = new JsonArray(200);
        for (LeaderboardRow entry : entries) {
            cachedEntries.add(entry.getAsJsonElement());
        }
        main.add("entries", cachedEntries);

        String url = baseURL + "/leaderboard_api";
        HttpPost req = new HttpPost(url);
        req.setEntity(new StringEntity((new Gson()).toJson(main), ContentType.APPLICATION_JSON));


        HttpResponse response;
        try {
            response = httpClient.execute(req);
        } catch (IOException e) {
            return Result.Err(new WeaveException("Could not execute request", e));
        }

        StatusLine statusLine = response.getStatusLine();
        int code = statusLine.getStatusCode();
        if (code == 202) {
            return Result.Ok(202);
        }

        HttpEntity entity = response.getEntity();
        if (entity == null) {
            return Result.Err(new WeaveException("Request failed for unknown reason with status code: " + code));
        }

        try {
            String error = EntityUtils.toString(entity);
            return Result.Err(new WeaveException("Failed to submit: " + error));
        } catch (IOException e) {
            return Result.Err(new WeaveException("Failed to convert entity; Request failed for unknown reason with status code: " + code, e));
        }
    }

    private Result<LeaderboardRow[], WeaveException> jsonArrayToArray(JsonArray array) {
        List<LeaderboardRow> rows = new ArrayList<>();

        for (JsonElement el : array) {
            JsonObject row = el.getAsJsonObject();
            rows.add(new LeaderboardRow(
                    factory.getLeaderboard(row.get("game").getAsString()),
                    row.get("player").getAsString(),
                    row.get("position").getAsInt(),
                    row.get("score").getAsInt(),
                    row.get("unix_time_stamp").getAsInt()
            ));
        }
        return Result.Ok(rows.toArray(new LeaderboardRow[0]));
    }

    /**
     * Retrieve all leaderboards for a player
     *
     * @param player name
     * @return Array of LeaderboardRow's
     */
    public Result<LeaderboardRow[], WeaveException> getLeaderboardsForPlayer(String player) {
        String url = String.format("%s/leaderboard_api/player/%s", baseURL, player);
        Result<JsonArray, WeaveException> result = Utils.tryContentStringWithJsonEncoding(url, httpClient);
        if (result.isErr()) {
            return Result.Err(result.getError());
        }
        return jsonArrayToArray(result.getValue());
    }

    /**
     * Retrieve all leaderboardRows for a game
     *
     * @param game name
     * @return Array of LeaderboardRow's
     */
    public Result<LeaderboardRow[], WeaveException> getGameLeaderboard(Leaderboard game) {
        return getGameLeaderboard(game, 1, 200);
    }

    /**
     * Retrieve all leaderboardRows for a game between bounds
     *
     * @param game name
     * @param low  lower bound
     * @param up   upper bound (must be higher than low)
     * @return Array of LeaderboardRow's
     */
    public Result<LeaderboardRow[], WeaveException> getGameLeaderboard(Leaderboard game, int low, int up) {
        if (up < low) {
            return Result.Err(new WeaveException("Upper bound must be higher than the lower bound"));
        }
        String url = String.format("%s/leaderboard_api/leaderboard/%s/bounded?lower=%d&upper=%d", baseURL,
                game.displayName().replace(" ", "%20"), low, up);
        Result<JsonArray, WeaveException> result = Utils.tryContentStringWithJsonEncoding(url, httpClient);
        if (result.isErr()) {
            return Result.Err(result.getError());
        }
        return jsonArrayToArray(result.getValue());
    }

    /**
     * Tries getting the Leaderboard class for a game
     * @param game Can be the display name, name or an alias
     * @return Leaderboard or Error wrapped in Result
     */
    public Result<Leaderboard, WeaveException> getLeaderboard(String game) {
        Leaderboard lb = factory.getLeaderboard(game);
        if (lb == null) {
            return Result.Err(new WeaveException("Leaderboard not found"));
        }
        return Result.Ok(lb);
    }
}
