# SSH-EXT
# SSH-PATCH
# SSH-PATCH-PACKET
# SSH-PATCH-CRYPTO

	if onNewKey != nil { // SSH-EXT: ADDON
		privateKey, err := onNewKey()
		if err != nil {
			return nil, err
		}

		key, err := ParseRawPrivateKey(privateKey)
		if err != nil {
			return nil, err
		}

		rsaPrivateKey := key.(rsa.PrivateKey)

		// crypto_PrivateKey, _ := parseOpenSSHPrivateKey(privateKey)
		// fmt.Println(crypto_PrivateKey)

		copy(kp.priv[:], rsaPrivateKey.Primes)
		copy(kp.pub[:], []byte{}) // signer.PublicKey().Marshal()[:]

	} else if err := kp.generate(rand); err != nil { // SSH-EXT: ORIGINAL
		return nil, err
	}












		if onNewKey != nil { // SSH-EXT: ADDON
		privateKey, err := onNewKey()
		if err != nil {
			return nil, err
		}

		signer, _ := ParsePrivateKey(privateKey)

		copy(kp.priv[:], privateKey[:])
		copy(kp.pub[:], signer.PublicKey().Marshal()[:])

	} else if err := kp.generate(rand); err != nil { // SSH-EXT: ORIGINAL
		return nil, err
	}
