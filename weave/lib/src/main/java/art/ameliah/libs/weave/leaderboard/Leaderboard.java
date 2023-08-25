package art.ameliah.libs.weave.leaderboard;

/**
 *
 * @param name internal name
 * @param displayName display name
 * @param active if the leaderboard is active (and can be submitted to)
 */
public record Leaderboard(String name, String displayName, boolean active, String scoreType) {}
