package art.ameliah.libs.weave;

import com.google.gson.Gson;
import com.google.gson.JsonArray;
import com.google.gson.JsonElement;
import com.google.gson.JsonObject;
import org.asynchttpclient.AsyncHttpClient;

import javax.annotation.Nullable;
import java.util.*;
import java.util.concurrent.CompletableFuture;

import static art.ameliah.libs.weave.Utils.makeRequest;

/**
 * Interaction class for LeaderboardAPI
 */
public class LeaderboardAPI {

    private final String baseURL;
    private final AsyncHttpClient httpClient;

    private final HashMap<String, Leaderboard> converter = new HashMap<>();

    /**
     * @param url        base API url
     * @param httpClient connection client
     */
    public LeaderboardAPI(String url, AsyncHttpClient httpClient) {
        this.baseURL = url;
        this.httpClient = httpClient;

        String url2 = String.format("%s/leaderboard_api/games/true", baseURL);
        CompletableFuture<JsonArray> completableFuture = makeRequest(httpClient, url2, JsonArray.class);
        completableFuture
                .whenComplete((leaderboards, throwable) -> {
                    if (throwable != null) {
                        return;
                    }

                    for (JsonElement el : leaderboards) {
                        JsonObject obj = el.getAsJsonObject();
                        Leaderboard game = new Leaderboard(obj.get("game").getAsString(), obj.get("display_name").getAsString(), true, obj.get("score_type").getAsString());
                        for (JsonElement aliasElement : obj.get("aliases").getAsJsonArray()) {
                            converter.put(aliasElement.getAsString(), game);
                        }
                        converter.put(game.name(), game);
                        converter.put(game.displayName(), game);
                    }
                });
    }

    /**
     * Submit a new leaderboard
     *
     * @param uuid    Submitters uuid
     * @param game    name
     * @param entries LeaderboardRows
     * @return StatusCode or Error
     */
    public CompletableFuture<Integer> submitLeaderboard(UUID uuid, Leaderboard game, Set<LeaderboardRow> entries) {
        if (entries.size() != 200) {
            return CompletableFuture.failedFuture(new WeaveException("entries set must contain exactly 200 entries."));
        }
        if (!game.active()) {
            return CompletableFuture.failedFuture(new WeaveException("Game is not active, cannot submit"));
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

        return httpClient.preparePost(url)
                .addHeader("Content-Type", "application/json")
                .setBody((new Gson()).toJson(main))
                .execute()
                .toCompletableFuture()
                .thenApplyAsync(response -> {
                    WeaveException e = WeaveException.fromResponse(response, 202, false);
                    if (e != null) {
                        throw new RuntimeException(e);
                    }

                    return 202;
                });
    }

    private LeaderboardRow[] jsonArrayToArray(JsonArray array) {
        List<LeaderboardRow> rows = new ArrayList<>();

        for (JsonElement el : array) {
            JsonObject row = el.getAsJsonObject();
            rows.add(new LeaderboardRow(
                    this.converter.get(row.get("game").getAsString()),
                    row.get("player").getAsString(),
                    row.get("position").getAsInt(),
                    row.get("score").getAsInt(),
                    row.get("unix_time_stamp").getAsInt()
            ));
        }
        return rows.toArray(new LeaderboardRow[0]);
    }

    private CompletableFuture<LeaderboardRow[]> leaderBoardRowRequest(String url) {
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
                    JsonArray array = (new Gson()).fromJson(json, JsonArray.class);
                    return jsonArrayToArray(array);
                });
    }

    /**
     * Retrieve all leaderboards for a player
     *
     * @param player name
     * @return Array of LeaderboardRow's
     */
    public CompletableFuture<LeaderboardRow[]> getLeaderboardsForPlayer(String player) {
        String url = String.format("%s/leaderboard_api/player/%s", baseURL, player);
        return leaderBoardRowRequest(url);
    }

    /**
     * Retrieve all leaderboardRows for a game
     *
     * @param game name
     * @return Array of LeaderboardRow's
     */
    public CompletableFuture<LeaderboardRow[]> getGameLeaderboard(Leaderboard game) {
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
    public CompletableFuture<LeaderboardRow[]> getGameLeaderboard(Leaderboard game, int low, int up) {
        if (up < low) {
            return CompletableFuture.failedFuture(new WeaveException("Upper bound must be higher than the lower bound"));
        }
        String url = String.format("%s/leaderboard_api/leaderboard/%s/bounded?lower=%d&upper=%d", baseURL,
                game.displayName().replace(" ", "%20"), low, up);
        return leaderBoardRowRequest(url);
    }

    /**
     * Tries getting the Leaderboard class for a game
     *
     * @param game Can be the display name, name or an alias
     * @return Leaderboard or Error wrapped in Result
     */
    public @Nullable Leaderboard getLeaderboard(String game) {
        return this.converter.get(game);
    }

    /**
     * @param name        internal name
     * @param displayName display name
     * @param active      if the leaderboard is active (and can be submitted to)
     */
    public record Leaderboard(String name, String displayName, boolean active, String scoreType) {
    }

    /**
     * Leaderboard Row
     *
     * @param game     Game's leaderboards
     * @param player   name
     * @param position int
     * @param score    int
     * @param unix     submission unix time stamp
     */
    public record LeaderboardRow(Leaderboard game, String player, int position, int score, int unix) {
        JsonElement getAsJsonElement() {
            JsonObject jsonObject = new JsonObject();
            jsonObject.addProperty("game", game.displayName());
            jsonObject.addProperty("player", player);
            jsonObject.addProperty("position", position);
            jsonObject.addProperty("score", score);
            return jsonObject;
        }
    }
}
