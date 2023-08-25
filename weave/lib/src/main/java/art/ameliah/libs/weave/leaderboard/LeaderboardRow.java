package art.ameliah.libs.weave.leaderboard;

import com.google.gson.JsonElement;
import com.google.gson.JsonObject;

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
