export function encrypt(plaintext) {

    let key = "starkcash";
    let encryptedText = "";
    
    for (let i = 0; i < plaintext.length; i += 2) {
        const plaintextByte = parseInt(plaintext.substr(i, 2), 16);
        const keyByte = parseInt(key.substr(i, 2), 16);
        const encryptedByte = plaintextByte ^ keyByte;
        encryptedText += ('00' + encryptedByte.toString(16)).slice(-2);
    }

    return encryptedText;
}

