package art.ameliah.libs.weave;

import art.ameliah.libs.weave.leaderboard.Leaderboard;
import art.ameliah.libs.weave.leaderboard.LeaderboardRow;
import org.junit.jupiter.api.Test;

import java.net.MalformedURLException;
import java.util.Arrays;
import java.util.List;
import java.util.function.BiFunction;

import static org.junit.jupiter.api.Assertions.*;

public class WeaveTest {


    private final Weave weave;

    public WeaveTest() throws WeaveException, MalformedURLException {
        if (System.getProperty("test.prod").equals("true")) {
            weave = Weave.Production();
        } else {
            weave = Weave.Dev(7070);
        }

    }

    @Test
    void testEggWarsMapAPI() {
        Result<EggWarsMapAPI.EggWarsMap, WeaveException> result = weave.getEggWarsMapAPI().getEggWarsMap("palace");
        assertTrue(result.isOk());

        Result<EggWarsMapAPI.EggWarsMap[], WeaveException> result2 = weave.getEggWarsMapAPI().getAllEggWarsMaps();
        assertTrue(result2.isOk());
        assertTrue(result2.getValue().length >= 13);
    }

    @Test
    void testChestAPISeasonGetters() {
        assertArrayEquals(
                weave.getChestAPI().getSeasons().unwrap_or_default(() -> new String[0]),
                new String[]{"migration_release", "birthday_event"});
        assertArrayEquals(
                weave.getChestAPI().getSeasons(ChestAPI.SeasonType.RUNNING)
                        .unwrap_or_default(() -> new String[0]),
                new String[]{"migration_release"});
    }


    private <T> boolean booleanSupplier(T[] result, T[] expected, BiFunction<T, List<T>, Boolean> checker) {
        for (T exp : expected) {
            if (!checker.apply(exp, List.of(result))) {
                fail(String.format("Could not find expected element (%s) in array (%s).", exp, Arrays.toString(result)));
                return false;
            }
        }
        return true;
    }

    @Test
    void testChestAPIChestLocationGetters() {
        ChestAPI.ChestLocation[] result = weave.getChestAPI().getCurrentChestLocations()
                .unwrap_or_default(() -> new ChestAPI.ChestLocation[0]);
        ChestAPI.ChestLocation[] expected = new ChestAPI.ChestLocation[]{
                new ChestAPI.ChestLocation("migration_release", 453, 54, -378),
                new ChestAPI.ChestLocation("migration_release", 417, 55, -547),
                new ChestAPI.ChestLocation("migration_release", 403, 54, -391),
                new ChestAPI.ChestLocation("migration_release", 494, 53, -329)
        };
        assertTrue(booleanSupplier(result, expected, (exp, results) -> results.contains(exp)));
    }

    @Test
    void testLeaderboardAPILeaderboardRowGameGetters() {
        Leaderboard ew = weave.getLeaderboardAPI().getLeaderboard("ew").unwrap();
        Leaderboard parkour = weave.getLeaderboardAPI().getLeaderboard("parkour").unwrap();
        LeaderboardRow[] result = weave.getLeaderboardAPI().getGameLeaderboard(ew)
                .unwrap_or_default(() -> new LeaderboardRow[0]);
        LeaderboardRow[] expected = new LeaderboardRow[]{
                new LeaderboardRow(ew, "Mivke", 1, 31000, 0),
                new LeaderboardRow(ew, "Fesa", 11, 0, 0),
        };
        assertTrue(booleanSupplier(result, expected, (exp, results) -> {
            for (LeaderboardRow res : results) {
                if (res.game().equals(exp.game())
                        && res.player().equals(exp.player())
                        && res.position() == exp.position()
                        && (res.score() == exp.score() || exp.score() == 0)
                        && (res.unix() == exp.unix() || exp.unix() == 0)) {
                    return true;
                }
            }
            return false;
        }));

        LeaderboardRow[] result2 = weave.getLeaderboardAPI().getLeaderboardsForPlayer("Mivke")
                .unwrap_or_default(() -> new LeaderboardRow[0]);
        LeaderboardRow[] expected2 = new LeaderboardRow[]{
                new LeaderboardRow(ew, "Mivke", 1, 31000, 0),
                new LeaderboardRow(parkour, "Mivke", 1, 2322, 0),
        };
        assertTrue(booleanSupplier(result2, expected2, (exp, results) -> {
            for (LeaderboardRow res : results) {
                if (res.game().name().equals(exp.game().name())
                        && res.player().equals(exp.player())
                        && res.position() == exp.position()
                        && (res.score() == exp.score() || exp.score() == 0)
                        && (res.unix() == exp.unix() || exp.unix() == 0)) {
                    return true;
                }
            }
            return false;
        }));
    }
}
