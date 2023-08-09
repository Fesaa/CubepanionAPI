package art.ameliah.libs.weave;


import com.google.gson.JsonArray;
import com.google.gson.JsonElement;
import com.google.gson.JsonObject;
import org.apache.http.impl.client.CloseableHttpClient;

import java.util.ArrayList;
import java.util.List;

import static art.ameliah.libs.weave.Utils.tryContentStringWithJsonEncoding;

public class EggWarsMapAPI {

    private final String baseURL;

    private final CloseableHttpClient httpClient;

    /**
     * @param url        base API url
     * @param httpClient connection client
     */
    protected EggWarsMapAPI(String url, CloseableHttpClient httpClient) {
        this.baseURL = url;
        this.httpClient = httpClient;
    }

    public Result<EggWarsMap[], WeaveException> getAllEggWarsMaps() {
        String url = baseURL + "/eggwars_map_api";
        Result<JsonArray, WeaveException> result = tryContentStringWithJsonEncoding(url, httpClient);
        if (result.isErr()) {
            return Result.Err(result.getError());
        }
        List<EggWarsMap> maps = new ArrayList<>();
        JsonArray array = result.getValue();
        for (JsonElement map : array) {
            maps.add(toEggWarsMap(map.getAsJsonObject()));
        }
        return Result.Ok(maps.toArray(new EggWarsMap[0]));
    }

    public Result<EggWarsMap, WeaveException> getEggWarsMap(String name) {
        String url = baseURL + "/eggwars_map_api/" + name;
        Result<JsonObject, WeaveException> result = tryContentStringWithJsonEncoding(url, httpClient, JsonObject.class);
        if (result.isErr()) {
            return Result.Err(result.getError());
        }
        return Result.Ok(toEggWarsMap(result.getValue()));
    }

    private EggWarsMap toEggWarsMap(JsonObject map) {
        List<Generator> generatorsList = new ArrayList<>();
        JsonArray generators = map.get("generators").getAsJsonArray();
        for (JsonElement gen : generators) {
            JsonObject generator = gen.getAsJsonObject();
            generatorsList.add(new Generator(
                    generator.get("unique_name").getAsString(),
                    generator.get("ordering").getAsInt(),
                    generator.get("gen_type").getAsString(),
                    generator.get("gen_location").getAsString(),
                    generator.get("level").getAsInt(),
                    generator.get("count").getAsInt()
            ));
        }
        return new EggWarsMap(
                map.get("unique_name").getAsString(),
                map.get("map_name").getAsString(),
                map.get("team_size").getAsInt(),
                map.get("build_limit").getAsInt(),
                map.get("colours").getAsString(),
                generatorsList.toArray(new Generator[0]));
    }


    public record EggWarsMap(String uniqueName, String mapName, int teamSize, int buildLimit, String colour, Generator[] generators) {};
    public record Generator(String uniqueName, int ordering, String genType, String location, int level, int count) {};
}
