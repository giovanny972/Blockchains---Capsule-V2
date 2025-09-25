export interface Coin {
  denom: string;
  amount: string;
}

export enum CapsuleType {
  UNKNOWN = 0,
  SAFE = 1,
  TIME_LOCK = 2,
  CONDITIONAL = 3,
  MULTI_SIG = 4,
  DEAD_MANS_SWITCH = 5,
}

export enum CapsuleStatus {
  UNKNOWN = 0,
  ACTIVE = 1,
  UNLOCKED = 2,
  EXPIRED = 3,
  CANCELLED = 4,
}

export interface Capsule {
  id: string;
  owner: string;
  capsule_type: CapsuleType;
  status: CapsuleStatus;
  title: string;
  description: string;
  data: Uint8Array;
  unlock_time?: Date;
  creation_time: Date;
  deposit: Coin[];
  unlock_conditions?: Uint8Array;
  required_signatures?: number;
  authorized_addresses?: string[];
  last_heartbeat?: Date;
  heartbeat_interval?: number;
  beneficiaries?: string[];
}

export interface KeyShare {
  capsule_id: string;
  share_id: number;
  holder: string;
  encrypted_share: Uint8Array;
  threshold: number;
  total_shares: number;
}

export interface ConditionContract {
  address: string;
  method: string;
  params: Uint8Array;
  expected_result: Uint8Array;
}

export interface Statistics {
  total_capsules: string;
  active_capsules: string;
  unlocked_capsules: string;
  expired_capsules: string;
  cancelled_capsules: string;
  total_value_locked: Coin[];
}

export interface MsgCreateCapsule {
  creator: string;
  capsule_type: CapsuleType;
  title: string;
  description: string;
  data: Uint8Array;
  unlock_time?: Date;
  deposit: Coin[];
  unlock_conditions?: Uint8Array;
  required_signatures?: number;
  authorized_addresses?: string[];
  heartbeat_interval?: number;
  beneficiaries?: string[];
}

export interface MsgUnlockCapsule {
  unlocker: string;
  capsule_id: string;
  key_shares?: KeyShare[];
  proof?: Uint8Array;
}