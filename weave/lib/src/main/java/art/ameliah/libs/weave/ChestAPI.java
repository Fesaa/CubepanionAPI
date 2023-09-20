package art.ameliah.libs.weave;

import org.asynchttpclient.AsyncHttpClient;

import java.util.concurrent.CompletableFuture;

import static art.ameliah.libs.weave.Utils.makeRequest;

/**
 * Interaction class for ChestAPI
 */
public class ChestAPI {

    private final String baseURL;
    private final AsyncHttpClient httpClient;


    /**
     * @param url        base API url
     * @param httpClient connection client
     */
    protected ChestAPI(String url, AsyncHttpClient httpClient) {
        this.baseURL = url;
        this.httpClient = httpClient;
    }

    /**
     * Requests chests locations for the current running season
     *
     * @return Array of ChestLocation's
     */
    public CompletableFuture<ChestLocation[]> getCurrentChestLocations() {
        String url = baseURL + "/chest_api/current";
        return makeRequest(httpClient, url, ChestLocation[].class);
    }

    /**
     * Get all chest locations for a specific season
     *
     * @param season The Season the request the chests for
     * @return Array of ChestLocation's
     */
    public CompletableFuture<ChestLocation[]> getChestLocationsForSeason(String season) {
        String url = baseURL + "/chest_api/season/" + season;
        return makeRequest(httpClient, url, ChestLocation[].class);
    }

    /**
     * Get all seasons
     *
     * @return Array of seasons (String)
     */
    public CompletableFuture<String[]> getSeasons() {
        return getSeasons(SeasonType.ALL);
    }

    /**
     * Get seasons bounded by the request type
     *
     * @param seasonType Request type
     * @return Array of seasons (String)
     */
    public CompletableFuture<String[]> getSeasons(SeasonType seasonType) {
        String url = baseURL + "/chest_api/seasons/" + seasonType.bool();
        return makeRequest(httpClient, url, String[].class);
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
     *
     * @param season_name Season the chest is connected to
     * @param x           x-coord
     * @param y           y-coord
     * @param z           z-coord
     */
    public record ChestLocation(String season_name, int x, int y, int z) {
    }
}
