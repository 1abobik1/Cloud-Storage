const DB_NAME = 'crypto-keys';
const STORE_NAME = 'keys';
const KEY_ID = 'aes-key';

export async function generateAndStoreKey(): Promise<void> {
    const key = await crypto.subtle.generateKey(
        {
            name: 'AES-GCM',
            length: 256,
        },
        true, // ✅ extractable = true
        ['encrypt', 'decrypt']
    );

    const jwk = await crypto.subtle.exportKey('jwk', key);

    const db = await openDatabase();
    const tx = db.transaction(STORE_NAME, 'readwrite');
    const store = tx.objectStore(STORE_NAME);
    await store.put(jwk, KEY_ID);
    db.close();
}

export async function getStoredKey(): Promise<CryptoKey> {
    return new Promise((resolve, reject) => {
        const request = indexedDB.open(DB_NAME, 1);

        request.onerror = () => reject(request.error);
        request.onsuccess = async () => {
            const db = request.result;
            const transaction = db.transaction(STORE_NAME, 'readonly');
            const store = transaction.objectStore(STORE_NAME);
            const getRequest = store.get(KEY_ID);

            getRequest.onerror = () => reject(getRequest.error);
            getRequest.onsuccess = async () => {
                const jwk = getRequest.result;
                if (!jwk) return reject(new Error('Ключ не найден'));

                try {
                    const key = await crypto.subtle.importKey(
                        'jwk',
                        jwk,
                        { name: 'AES-GCM' },
                        true,
                        ['encrypt', 'decrypt']
                    );
                    resolve(key);
                } catch (e) {
                    reject(e);
                }
            };
        };

        request.onupgradeneeded = () => {
            const db = request.result;
            db.createObjectStore(STORE_NAME);
        };
    });
}

async function openDatabase(): Promise<IDBDatabase> {
    return new Promise((resolve, reject) => {
        const request = indexedDB.open(DB_NAME, 1);

        request.onupgradeneeded = () => {
            const db = request.result;
            db.createObjectStore(STORE_NAME);
        };

        request.onsuccess = () => resolve(request.result);
        request.onerror = () => reject(request.error);
    });
}
