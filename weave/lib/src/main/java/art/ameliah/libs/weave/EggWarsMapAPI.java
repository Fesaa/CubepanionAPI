package art.ameliah.libs.weave;


import org.asynchttpclient.AsyncHttpClient;

import java.util.concurrent.CompletableFuture;

import static art.ameliah.libs.weave.Utils.makeRequest;

/**
 * API for EggWars maps
 */
public class EggWarsMapAPI {

    private final String baseURL;

    private final AsyncHttpClient httpClient;

    /**
     * @param url        base API url
     * @param httpClient connection client
     */
    protected EggWarsMapAPI(String url, AsyncHttpClient httpClient) {
        this.baseURL = url;
        this.httpClient = httpClient;
    }

    /**
     * Get all EggWars maps
     *
     * @return all EggWars maps
     */
    public CompletableFuture<EggWarsMap[]> getAllEggWarsMaps() {
        String url = baseURL + "/eggwars_map_api";
        return makeRequest(httpClient, url, EggWarsMap[].class);
    }

    /**
     * Get EggWars map by name
     *
     * @param name map name
     * @return EggWars map
     */
    public CompletableFuture<EggWarsMap> getEggWarsMap(String name) {
        String url = baseURL + "/eggwars_map_api/" + name;
        return makeRequest(httpClient, url, EggWarsMap.class);
    }

    /**
     * EggWars map
     */
    public record EggWarsMap(String unique_name, String map_name, String layout, int team_size, int build_limit,
                             String colour, Generator[] generators) {
    }

    /**
     * Generator
     */
    public record Generator(String unique_name, int ordering, String gen_type, String gen_location, int level,
                            int count) {
    }

}
