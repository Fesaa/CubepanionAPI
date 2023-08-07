package art.ameliah.libs.weave;

import org.junit.jupiter.api.Test;
import art.ameliah.libs.weave.LeaderboardAPI.Leaderboard;

import java.net.MalformedURLException;
import java.util.List;
import java.util.function.BiFunction;

import static org.junit.jupiter.api.Assertions.*;
import static org.junit.jupiter.api.Assertions.assertTrue;

public class WeaveTest {


    private final Weave weave;

    public WeaveTest() throws MalformedURLException {
        if (System.getProperty("test.prod").equals("true")) {
            weave = Weave.Production();
        } else {
            weave = Weave.Dev();
        }

    }

    @Test
    void testChestAPISeasonGetters() throws WeaveException {
        assertArrayEquals(weave.getChestAPI().getSeasons(), new String[]{"migration_release", "birthday_event"});
        assertArrayEquals(weave.getChestAPI().getSeasons(ChestAPI.SeasonType.RUNNING), new String[]{"migration_release"});
    }


    private <T> boolean booleanSupplier(T[] result, T[] expected, BiFunction<T, List<T>, Boolean> checker) {
        for (T exp : expected) {
            if (!checker.apply(exp, List.of(result))) {
                fail(String.format("Could not find expected element (%s) in array.", exp));
                return false;
            }
        }
        return true;
    }

    @Test
    void testChestAPIChestLocationGetters() throws WeaveException {
        ChestAPI.ChestLocation[] result = weave.getChestAPI().getCurrentChestLocations();
        ChestAPI.ChestLocation[] expected = new ChestAPI.ChestLocation[]{
                new ChestAPI.ChestLocation("migration_release", 453, 54, -378),
                new ChestAPI.ChestLocation("migration_release", 417, 55, -547),
                new ChestAPI.ChestLocation("migration_release", 403, 54, -391),
                new ChestAPI.ChestLocation("migration_release", 494, 53, -329)
        };
        assertTrue(booleanSupplier(result, expected, (exp, results) -> results.contains(exp)));
    }

    @Test
    void testLeaderboardAPILeaderboardRowGameGetters() throws WeaveException {
        LeaderboardAPI.LeaderboardRow[] result = weave.getLeaderboardAPI().getGameLeaderboard(Leaderboard.TEAM_EGGWARS);
        LeaderboardAPI.LeaderboardRow[] expected = new LeaderboardAPI.LeaderboardRow[]{
                new LeaderboardAPI.LeaderboardRow(Leaderboard.TEAM_EGGWARS, "Mivke", 1, 31000, 0),
                new LeaderboardAPI.LeaderboardRow(Leaderboard.TEAM_EGGWARS, "Fesa", 11, 0, 0),
        };
        assertTrue(booleanSupplier(result, expected, (exp, results) -> {
            for (LeaderboardAPI.LeaderboardRow res: results) {
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

        LeaderboardAPI.LeaderboardRow[] result2 = weave.getLeaderboardAPI().getLeaderboardsForPlayer("Mivke");
        LeaderboardAPI.LeaderboardRow[] expected2 = new LeaderboardAPI.LeaderboardRow[]{
                new LeaderboardAPI.LeaderboardRow(Leaderboard.TEAM_EGGWARS, "Mivke", 1, 31000, 0),
                new LeaderboardAPI.LeaderboardRow(Leaderboard.PARKOUR, "Mivke", 12, 2277 , 0),
        };
        assertTrue(booleanSupplier(result2, expected2, (exp, results) -> {
            for (LeaderboardAPI.LeaderboardRow res: results) {
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
    }
}
