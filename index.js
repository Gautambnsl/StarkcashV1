export function encrypt(plaintext) {
    let key = 3;
    let encryptedText = "";
    
    for (let i = 0; i < plaintext.length; i += 2) {
        const plaintextByte = parseInt(plaintext.substr(i, 2), 16);
        const keyByte = parseInt(key.substr(i, 2), 16);
        const encryptedByte = plaintextByte ^ keyByte;
        encryptedText += encryptedByte.toString(16).padStart(2, '0');
    }

    return encryptedText;
}