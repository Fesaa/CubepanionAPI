package art.ameliah.libs.weave.leaderboard;

import javax.annotation.Nullable;

public class LeaderboardFactoryDummy implements ILeaderboardFactory{
    @Nullable
    @Override
    public Leaderboard getLeaderboard(String name) {
        return null;
    }
}
