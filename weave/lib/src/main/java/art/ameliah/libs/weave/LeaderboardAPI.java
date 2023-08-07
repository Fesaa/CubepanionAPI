package art.ameliah.libs.weave;

import com.google.gson.*;
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

    /**
     *
     * @param url base API url
     * @param httpClient connection client
     */
    protected LeaderboardAPI(String url, CloseableHttpClient httpClient) {
        this.baseURL = url;
        this.httpClient = httpClient;
    }

    /**
     * Submit a new leaderboard
     * @param uuid Submitters uuid
     * @param game name
     * @param entries LeaderboardRows
     * @throws WeaveException Any exceptions throw during requests
     */
    public void submitLeaderboard(UUID uuid, Leaderboard game, Set<LeaderboardRow> entries) throws WeaveException {
        if (entries.size() != 200) {
            throw new WeaveException("entries set must contain exactly 200 entries.");
        }

        JsonObject main = new JsonObject();
        main.addProperty("uuid", uuid.toString());
        main.addProperty("unix_time_stamp", System.currentTimeMillis());
        main.addProperty("game", game.getString());
        JsonArray cachedEntries = new JsonArray(200);
        for (LeaderboardRow entry : entries) {
            cachedEntries.add(entry.getAsJsonElement());
        }
        main.add("entries", cachedEntries);

        String url = baseURL + "/leaderboard_api";
        HttpPost req = new HttpPost(url);
        req.setEntity(new StringEntity(main.getAsString(), ContentType.APPLICATION_JSON));


        HttpResponse response;
        try {
            response = httpClient.execute(req);
        } catch (IOException e) {
            throw new WeaveException("Could not execute request", e);
        }

        StatusLine statusLine = response.getStatusLine();
        int code = statusLine.getStatusCode();
        if (code == 200) {
            return;
        }

        HttpEntity entity = response.getEntity();
        if (entity == null) {
            throw new WeaveException("Request failed for unknown reason with status code: " + code);
        }

        try {
            String error = EntityUtils.toString(entity);
            throw new WeaveException("Failed to submit: " + error);
        } catch (IOException e) {
            throw new WeaveException("Failed to convert entity; Request failed for unknown reason with status code: " + code, e);
        }
    }

    private LeaderboardRow[] jsonArrayToArray(JsonArray array) throws WeaveException {
        List<LeaderboardRow> rows = new ArrayList<>();

        for(JsonElement el : array) {
            JsonObject row = el.getAsJsonObject();
            rows.add(new LeaderboardRow(
                    Leaderboard.stringToLeaderboard(row.get("game").getAsString()),
                    row.get("player").getAsString(),
                    row.get("position").getAsInt(),
                    row.get("score").getAsInt(),
                    row.get("unix_time_stamp").getAsInt()
            ));
        }
        return rows.toArray(new LeaderboardRow[0]);
    }

    /**
     * Retrieve all leaderboards for a player
     * @param player name
     * @return Array of LeaderboardRow's
     * @throws WeaveException Any exceptions throw during requests
     */
    public LeaderboardRow[] getLeaderboardsForPlayer(String player) throws WeaveException {
        String url = String.format("%s/leaderboard_api/player/%s", baseURL, player);
        JsonArray array = Utils.tryContentStringWithJsonEncoding(url, httpClient);
        return jsonArrayToArray(array);
    }

    /**
     * Retrieve all leaderboardRows for a game
     * @param game name
     * @return Array of LeaderboardRow's
     * @throws WeaveException Any exceptions throw during requests
     */
    public LeaderboardRow[] getGameLeaderboard(Leaderboard game) throws WeaveException {
        return getGameLeaderboard(game, 1, 200);
    }

    /**
     * Retrieve all leaderboardRows for a game between bounds
     * @param game name
     * @param low lower bound
     * @param up upper bound (must be higher than low)
     * @return Array of LeaderboardRow's
     * @throws WeaveException Any exceptions throw during requests
     */
    public LeaderboardRow[] getGameLeaderboard(Leaderboard game, int low, int up) throws WeaveException {
        if (up < low) throw new WeaveException("Upper bound must be higher than the lower bound");
        String url = String.format("%s/leaderboard_api/leaderboard/%s/bounded?lower=%d&upper=%d", baseURL,
                game.getString().replace(" ", "%20"), low, up);
        JsonArray array = Utils.tryContentStringWithJsonEncoding(url, httpClient);
        return jsonArrayToArray(array);
    }

    /**
     * Leaderboard Row
     * @param game Game's leaderboards
     * @param player name
     * @param position int
     * @param score int
     * @param unix submission unix time stamp
     */
    public record LeaderboardRow(Leaderboard game, String player, int position, int score, int unix) {
        private JsonElement getAsJsonElement() {
            JsonObject jsonObject = new JsonObject();
            jsonObject.addProperty("game", game.getString());
            jsonObject.addProperty("player", player);
            jsonObject.addProperty("position", position);
            jsonObject.addProperty("score", score);
            return jsonObject;
        }
    }

    /**
     * Avaiable leaderboards
     */
    public enum Leaderboard {
        /**
         * In use
         */
        TEAM_EGGWARS("Team EggWars"),
        /**
         * Not updated. Will return no rows
         */
        TEAM_EGGWARS_SEASON_2("Team EggWars Season 2"),
        /**
         * In use
         */
        SOLO_LUCKYISLANDS("Lucky Islands"),
        /**
         * In use
         */
        SOLO_SKYWARS("Solo SkyWars"),
        /**
         * In use
         */
        FFA("Free For All"),
        /**
         * For completion
         */
        NONE(""),
        /**
         * In use
         */
        PARKOUR("Parkour");


        private final String string;

        Leaderboard(String s) {
            this.string = s;
        }

        /**
         /**
         * Tries converting to a leaderboard
         * @param s String to try on
         * @return Leaderboard (NONE if none found)
         * @throws WeaveException If no Leaderboard are found
         */
        public static Leaderboard stringToLeaderboard(String s) throws WeaveException {
            switch (s.toLowerCase()) {
                case "team eggwars", "eggwars", "tew", "ew" -> {
                    return Leaderboard.TEAM_EGGWARS;
                }
                case "eggwars season 2", "ew2", "tew2" -> {
                    return Leaderboard.TEAM_EGGWARS_SEASON_2;
                }
                case "solo skywars", "skywars", "sw" -> {
                    return Leaderboard.SOLO_SKYWARS;
                }
                case "lucky islands", "li" -> {
                    return Leaderboard.SOLO_LUCKYISLANDS;
                }
                case "free for all", "ffa" -> {
                    return Leaderboard.FFA;
                }
                case "parkour" -> {
                    return Leaderboard.PARKOUR;
                }
                default -> {
                    throw new WeaveException("String is not a valid leaderboard");
                }
            }
        }

        /**
         *
         * @return Properly formatted String
         */
        public String getString() {
            return string;
        }
    }

}
