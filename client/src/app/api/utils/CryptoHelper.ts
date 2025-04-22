const SECRET_KEY_STRING = process.env.SECRET_KEY!;
const ENCODER = new TextEncoder();

new TextDecoder();
async function getCryptoKey(): Promise<CryptoKey> {
    const keyData = ENCODER.encode(SECRET_KEY_STRING.padEnd(32, '0').slice(0, 32)); // 256-bit key
    return await crypto.subtle.importKey(
        'raw',
        keyData,
        'AES-GCM',
        false,
        ['encrypt', 'decrypt']
    );
}

function generateIV(): Uint8Array {
    return crypto.getRandomValues(new Uint8Array(12)); // 96-bit IV
}

export const cryptoHelper = {
    async encryptFile(file: File): Promise<File> {
        const arrayBuffer = await file.arrayBuffer();
        const iv = generateIV();
        const key = await getCryptoKey();

        const encrypted = await crypto.subtle.encrypt(
            { name: 'AES-GCM', iv },
            key,
            arrayBuffer
        );

        // комбинируем IV + зашифрованные данные
        const combined = new Uint8Array(iv.length + encrypted.byteLength);
        combined.set(iv, 0);
        combined.set(new Uint8Array(encrypted), iv.length);

        return new File([combined], file.name, { type: file.type });
    },

    async decryptFile(file: File): Promise<Blob> {
        const combined = new Uint8Array(await file.arrayBuffer());
        const iv = combined.slice(0, 12);
        const data = combined.slice(12);

        const key = await getCryptoKey();

        const decrypted = await crypto.subtle.decrypt(
            { name: 'AES-GCM', iv },
            key,
            data
        );

        // Важно: возвращаем Blob с оригинальным MIME type
        return new Blob([decrypted], { type: file.type });
    }

};
