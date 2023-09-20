package art.ameliah.libs.weave;

/**
 * Special WeaveException for Rate Limits
 */
public class RateLimited extends WeaveException {
    /**
     * RateLimited constructor
     *
     * @param seconds time out (seconds as this is returned)
     */
    public RateLimited(String seconds) {
        super(String.format("You have been rate limited, try again in %s seconds.", seconds));
    }
}
