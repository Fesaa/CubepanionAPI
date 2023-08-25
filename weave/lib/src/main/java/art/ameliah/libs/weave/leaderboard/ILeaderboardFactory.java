package art.ameliah.libs.weave.leaderboard;

import javax.annotation.Nullable;

public interface ILeaderboardFactory {

    @Nullable
    Leaderboard getLeaderboard(String name);
}
