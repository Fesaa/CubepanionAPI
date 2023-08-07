package art.ameliah.libs.weave;

import com.google.gson.*;
import org.apache.http.impl.client.CloseableHttpClient;

import java.util.ArrayList;
import java.util.List;

/**
 * Interaction class for ChestAPI
 */
public class ChestAPI {

    private final String baseURL;
    private final CloseableHttpClient httpClient;


    /**
     *
     * @param url base API url
     * @param httpClient connection client
     */
    protected ChestAPI(String url, CloseableHttpClient httpClient) {
        this.baseURL = url;
        this.httpClient = httpClient;
    }

    /**
     * Requests chests locations for the current running season
     * @return Array of ChestLocation's
     * @throws WeaveException Any exceptions throw during requests
     */
    public ChestLocation[] getCurrentChestLocations() throws WeaveException {
        JsonArray array = Utils.tryContentStringWithJsonEncoding(baseURL + "/chest_api/current", httpClient);
        List<ChestLocation> locs = new ArrayList<>();

        for (JsonElement el : array) {
            JsonObject chestLocation = el.getAsJsonObject();
            locs.add(new ChestLocation(
                    chestLocation.get("season_name").getAsString(),
                    chestLocation.get("x").getAsInt(),
                    chestLocation.get("y").getAsInt(),
                    chestLocation.get("z").getAsInt()
            ));
        }
        return locs.toArray(new ChestLocation[0]);
    }

    /**
     * Get all chest locations for a specific season
     * @param season The Season the request the chests for
     * @return Array of ChestLocation's
     * @throws WeaveException Any exceptions throw during requests
     */
    public ChestLocation[] getChestLocationsForSeason(String season) throws WeaveException {
        JsonArray array = Utils.tryContentStringWithJsonEncoding(baseURL + "/chest_api/season/" + season, httpClient);
        try {
            return new Gson().fromJson(array, ChestLocation[].class);
        } catch (JsonSyntaxException e) {
            throw new WeaveException("Could not constrict ChestLocation[]", e);
        }
    }

    /**
     * Get all seasons
     * @return Array of seasons (String)
     * @throws WeaveException Any exceptions throw during requests
     */
    public String[] getSeasons() throws WeaveException {
        return getSeasons(SeasonType.ALL);
    }

    /**
     * Get seasons bounded by the request type
     * @param seasonType Request type
     * @return Array of seasons (String)
     * @throws WeaveException Any exceptions throw during requests
     */
    public String[] getSeasons(SeasonType seasonType) throws WeaveException {
        JsonArray array = Utils.tryContentStringWithJsonEncoding(baseURL + "/chest_api/seasons/" + seasonType.bool(), httpClient);
        try {
            return new Gson().fromJson(array, String[].class);
        } catch (JsonSyntaxException e) {
            throw new WeaveException("Could not constrict ChestLocation[]", e);
        }
    }

    /**
     * Request type for seasons
     */
    public enum SeasonType {
        /**
         * Active season
         */
        RUNNING,
        /**
         * All seasons
         */
        ALL;

        private String bool() {
            if (this == SeasonType.RUNNING) {
                return "true";
            }
            return "false";
        }
    }

    /**
     * Location of a chest
     * @param season_name Season the chest is connected to
     * @param x x-coord
     * @param y y-coord
     * @param z z-coord
     */
    public record ChestLocation(String season_name, int x, int y, int z) {
    }
}
