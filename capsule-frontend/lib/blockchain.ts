import { StargateClient, SigningStargateClient } from "@cosmjs/stargate";
import { DirectSecp256k1HdWallet, OfflineDirectSigner } from "@cosmjs/proto-signing";
import { stringToPath } from "@cosmjs/crypto";
import { Window as KeplrWindow } from "@keplr-wallet/types";

declare global {
  interface Window extends KeplrWindow {}
}

export const CHAIN_ID = "capsule-mainnet";
export const RPC_ENDPOINT = "http://localhost:26657";
export const REST_ENDPOINT = "http://localhost:1317";
export const BECH32_PREFIX = "cosmos";

export const CHAIN_INFO = {
  chainId: CHAIN_ID,
  chainName: "Capsule Network",
  rpc: RPC_ENDPOINT,
  rest: REST_ENDPOINT,
  bip44: {
    coinType: 118,
  },
  bech32Config: {
    bech32PrefixAccAddr: BECH32_PREFIX,
    bech32PrefixAccPub: BECH32_PREFIX + "pub",
    bech32PrefixValAddr: BECH32_PREFIX + "valoper",
    bech32PrefixValPub: BECH32_PREFIX + "valoperpub",
    bech32PrefixConsAddr: BECH32_PREFIX + "valcons",
    bech32PrefixConsPub: BECH32_PREFIX + "valconspub",
  },
  currencies: [
    {
      coinDenom: "CAPS",
      coinMinimalDenom: "stake",
      coinDecimals: 6,
    },
  ],
  feeCurrencies: [
    {
      coinDenom: "CAPS",
      coinMinimalDenom: "stake",
      coinDecimals: 6,
      gasPriceStep: {
        low: 0.01,
        average: 0.025,
        high: 0.04,
      },
    },
  ],
  stakeCurrency: {
    coinDenom: "CAPS",
    coinMinimalDenom: "stake",
    coinDecimals: 6,
  },
};

export class BlockchainClient {
  private client: StargateClient | null = null;
  private signingClient: SigningStargateClient | null = null;

  async connect(): Promise<StargateClient> {
    if (!this.client) {
      this.client = await StargateClient.connect(RPC_ENDPOINT);
    }
    return this.client;
  }

  async connectWithKeplr(): Promise<{
    client: SigningStargateClient;
    address: string;
  }> {
    if (typeof window === "undefined" || !window.keplr) {
      throw new Error("Keplr wallet not found");
    }

    try {
      // Suggest chain to Keplr
      await window.keplr.experimentalSuggestChain(CHAIN_INFO);
    } catch (error) {
      console.warn("Failed to suggest chain to Keplr:", error);
    }

    // Enable chain
    await window.keplr.enable(CHAIN_ID);

    // Get offline signer
    const offlineSigner = window.keplr.getOfflineSigner(CHAIN_ID);
    const accounts = await offlineSigner.getAccounts();

    if (accounts.length === 0) {
      throw new Error("No accounts found in Keplr");
    }

    // Create signing client
    const client = await SigningStargateClient.connectWithSigner(
      RPC_ENDPOINT,
      offlineSigner,
      {
        gasPrice: {
          denom: "stake",
          amount: "0.025",
        },
      }
    );

    this.signingClient = client;

    return {
      client,
      address: accounts[0].address,
    };
  }

  async createWalletFromMnemonic(mnemonic: string): Promise<{
    client: SigningStargateClient;
    address: string;
  }> {
    const wallet = await DirectSecp256k1HdWallet.fromMnemonic(mnemonic, {
      prefix: BECH32_PREFIX,
      hdPaths: [stringToPath("m/44'/118'/0'/0/0")],
    });

    const [account] = await wallet.getAccounts();
    const client = await SigningStargateClient.connectWithSigner(
      RPC_ENDPOINT,
      wallet,
      {
        gasPrice: {
          denom: "stake",
          amount: "0.025",
        },
      }
    );

    this.signingClient = client;

    return {
      client,
      address: account.address,
    };
  }

  async getBalance(address: string): Promise<any> {
    const client = await this.connect();
    return await client.getAllBalances(address);
  }

  async getChainInfo(): Promise<any> {
    const client = await this.connect();
    return {
      chainId: await client.getChainId(),
      height: await client.getHeight(),
    };
  }

  // REST API methods
  async queryCapsules(): Promise<any> {
    const response = await fetch(`${REST_ENDPOINT}/cosmos/timecapsule/v1/capsules`);
    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`);
    }
    return await response.json();
  }

  async queryCapsule(id: string): Promise<any> {
    const response = await fetch(`${REST_ENDPOINT}/cosmos/timecapsule/v1/capsules/${id}`);
    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`);
    }
    return await response.json();
  }

  async queryUserCapsules(address: string): Promise<any> {
    const response = await fetch(`${REST_ENDPOINT}/cosmos/timecapsule/v1/capsules/user/${address}`);
    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`);
    }
    return await response.json();
  }

  async queryStats(): Promise<any> {
    const response = await fetch(`${REST_ENDPOINT}/cosmos/timecapsule/v1/stats`);
    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`);
    }
    return await response.json();
  }

  // Transaction methods
  async createCapsule(
    senderAddress: string,
    message: any,
    fee: any = "auto",
    memo: string = ""
  ): Promise<any> {
    if (!this.signingClient) {
      throw new Error("Signing client not connected");
    }

    const msgCreateCapsule = {
      typeUrl: "/cosmos.timecapsule.v1.MsgCreateCapsule",
      value: message,
    };

    return await this.signingClient.signAndBroadcast(
      senderAddress,
      [msgCreateCapsule],
      fee,
      memo
    );
  }

  async unlockCapsule(
    senderAddress: string,
    capsuleId: string,
    fee: any = "auto",
    memo: string = ""
  ): Promise<any> {
    if (!this.signingClient) {
      throw new Error("Signing client not connected");
    }

    const msgUnlockCapsule = {
      typeUrl: "/cosmos.timecapsule.v1.MsgUnlockCapsule",
      value: {
        unlocker: senderAddress,
        capsuleId: capsuleId,
      },
    };

    return await this.signingClient.signAndBroadcast(
      senderAddress,
      [msgUnlockCapsule],
      fee,
      memo
    );
  }

  disconnect() {
    if (this.client) {
      this.client.disconnect();
      this.client = null;
    }
    if (this.signingClient) {
      this.signingClient.disconnect();
      this.signingClient = null;
    }
  }
}