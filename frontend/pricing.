interface PriceResponse {
    price: number;
    currency?: string;
}

interface CacheEntry {
    price: number;
    timestamp: number;
}

class PricingError extends Error {
    constructor(message: string, public readonly productId: string, public readonly cause?: unknown) {
        super(message);
        this.name = 'PricingError';
    }
}

export class PricingService {
    private cache = new Map<string, CacheEntry>();
    private inFlight = new Map<string, Promise<number>>();
    private failureCount = 0;
    private circuitOpenUntil = 0;

    private readonly CACHE_TTL_MS = 30_000;
    private readonly CIRCUIT_THRESHOLD = 5;
    private readonly CIRCUIT_COOLDOWN_MS = 15_000;
    private readonly MAX_RETRIES = 3;
    private readonly TIMEOUT_MS = 8_000;

    async fetchDynamicPrice(productId: string): Promise<number> {
        // 1. Проверка кэша
        const cached = this.cache.get(productId);
        if (cached && Date.now() - cached.timestamp < this.CACHE_TTL_MS) {
            return cached.price;
        }

        // 2. Circuit breaker
        if (Date.now() < this.circuitOpenUntil) {
            throw new PricingError('Circuit open — service temporarily unavailable', productId);
        }

        // 3. Дедупликация параллельных запросов на один и тот же productId
        const existing = this.inFlight.get(productId);
        if (existing) return existing;

        const promise = this.fetchWithRetry(productId)
            .then(price => {
                this.cache.set(productId, { price, timestamp: Date.now() });
                this.failureCount = 0;
                return price;
            })
            .catch(err => {
                this.failureCount++;
                if (this.failureCount >= this.CIRCUIT_THRESHOLD) {
                    this.circuitOpenUntil = Date.now() + this.CIRCUIT_COOLDOWN_MS;
                }
                throw err;
            })
            .finally(() => {
                this.inFlight.delete(productId);
            });

        this.inFlight.set(productId, promise);
        return promise;
    }

    private async fetchWithRetry(productId: string, attempt = 1): Promise<number> {
        try {
            return await this.fetchOnce(productId);
        } catch (err) {
            if (attempt >= this.MAX_RETRIES) {
                throw new PricingError(`Failed after ${attempt} attempts`, productId, err);
            }
            const backoff = Math.min(1000 * 2 ** (attempt - 1), 4000) + Math.random() * 300;
            await new Promise(r => setTimeout(r, backoff));
            return this.fetchWithRetry(productId, attempt + 1);
        }
    }

    private async fetchOnce(productId: string): Promise<number> {
        const controller = new AbortController();
        const timeout = setTimeout(() => controller.abort(), this.TIMEOUT_MS);

        try {
            const response = await fetch(`/api/v1/products/${productId}/price`, {
                method: 'GET',
                headers: {
                    'Content-Type': 'application/json',
                    'X-Client-ID': 'zuhro-frontend'
                },
                signal: controller.signal
            });

            if (!response.ok) {
                throw new PricingError(`HTTP ${response.status}`, productId);
            }

            const data: PriceResponse = await response.json();

            if (typeof data.price !== 'number' || isNaN(data.price) || data.price < 0) {
                throw new PricingError('Invalid price payload', productId);
            }

            return data.price;
        } catch (err) {
            if (err instanceof Error && err.name === 'AbortError') {
                throw new PricingError('Request timeout', productId, err);
            }
            throw err;
        } finally {
            clearTimeout(timeout);
        }
    }

    clearCache(productId?: string): void {
        productId ? this.cache.delete(productId) : this.cache.clear();
    }
}
