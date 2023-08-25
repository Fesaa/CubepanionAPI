package art.ameliah.libs.weave;

import art.ameliah.libs.weave.leaderboard.LeaderboardAPI;
import org.apache.http.impl.client.CloseableHttpClient;
import org.apache.http.impl.client.HttpClients;

import java.net.MalformedURLException;
import java.net.URL;

/**
 * API access class
 */
public class Weave {

    private final LeaderboardAPI leaderboardAPI;
    private final ChestAPI chestAPI;
    private final EggWarsMapAPI eggWarsMapAPI;

    private Weave(String domain, int port) throws MalformedURLException {
        String baseURL = (new URL(String.format("http://%s:%d", domain, port))).toString();
        CloseableHttpClient httpClient = HttpClients.createDefault();
        leaderboardAPI = new LeaderboardAPI(baseURL, httpClient);
        chestAPI = new ChestAPI(baseURL, httpClient);
        eggWarsMapAPI = new EggWarsMapAPI(baseURL, httpClient);
    }

    /**
     * Setup API in prod env
     * @return Weave
     */
    public static Weave Production() {
        try {
            return new Weave("ameliah.art", 7070);
        } catch (MalformedURLException e) {
            throw new RuntimeException(e);
        }
    }

    /**
     * Setup API for the default dev env (http://127.0.0.1:8080)
     *
     * @return Weave
     */
    public static Weave Dev() {
        try {
            return new Weave("127.0.0.1", 8080);
        } catch (MalformedURLException e) {
            throw new RuntimeException(e);
        }
    }

    /**
     * Setup API in dev env with custom port
     *
     * @param port custom port
     * @return Weave
     * @throws MalformedURLException Could not construct API-url
     */
    public static Weave Dev(int port) throws MalformedURLException {
        return new Weave("127.0.0.1", port);
    }

    /**
     * Setup custom API in dev env
     *
     * @param domain custom domain
     * @param port   custom port
     * @return Weave
     * @throws MalformedURLException Could not construct API-url
     */
    public static Weave Dev(String domain, int port) throws MalformedURLException {
        return new Weave(domain, port);
    }

    /**
     * @return LeaderboardAPI
     */
    public LeaderboardAPI getLeaderboardAPI() {
        return leaderboardAPI;
    }

    /**
     * @return ChestAPI
     */
    public ChestAPI getChestAPI() {
        return chestAPI;
    }

    /**
     * @return EggWarsMapAPI
     */
    public EggWarsMapAPI getEggWarsMapAPI() {
        return eggWarsMapAPI;
    }
}
