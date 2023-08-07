package art.ameliah.libs.weave;

import javax.annotation.Nullable;
import java.util.function.Supplier;

/**
 * Rust like result class to handle error nicer
 * @param <T> Possible value
 * @param <E> Possible error
 */
public class Result<T, E> {

    @Nullable
    private final T value;

    @Nullable
    private final E error;

    private boolean checked = false;

    /**
     * Constructor
     * @param v value
     * @param e error
     */
    protected Result(@Nullable T v, @Nullable E e) {
        value = v;
        error = e;
    }

    /**
     * Construct Result containing value
     * @param v value
     * @return the result
     * @param <T> result type
     * @param <E> error type (absent)
     */
    public static <T, E> Result<T, E> Ok(T v) {
        return new Result<>(v, null);
    }

    /**
     * Construct Result containing error
     * @param e error
     * @return the result
     * @param <T> result type (absent)
     * @param <E> error type
     */
    public static <T, E> Result<T, E> Err(E e) {
        return new Result<>(null, e);
    }

    /**
     *
     * @return If it's Err
     */
    public boolean isErr() {
        checked = true;
        return value == null;
    }

    /**
     *
     * @return If it's Ok
     */
    public boolean isOk() {
        checked = true;
        return error == null;
    }

    /**
     * Will throw a RuntimeException if the Result hasn't been checked or contains a value
     * @return The error
     */
    public E getError() {
        if (!checked) {
            throw new RuntimeException("Tried accessing error before checking");
        }
        if (value != null) {
            throw new RuntimeException("Cannot access error, value present");
        }
        return error;
    }

    /**
     * Will throw a RuntimeException if the Result hasn't been checked or contains an error
     * @return The value
     */
    public T getValue() {
        if (!checked) {
            throw new RuntimeException("Tried accessing value before checking");
        }
        if (error != null) {
            throw new RuntimeException("Cannot access value, error present");
        }
        return value;
    }

    /**
     * Unsafely get the value, RuntimeException if error present
     * @return value
     */
    public T unwrap() {
        if (error != null) {
            throw new RuntimeException("Cannot access value, error present");
        }
        return value;
    }

    /**
     * Safely get the value
     * @param defaultSupplier value if error is present
     * @return value
     */
    public T unwrap_or_default(Supplier<T> defaultSupplier) {
        if (error != null) {
            return defaultSupplier.get();
        }
        return value;
    }

}
