import { MeetingInfo, MeetingUser } from '@/models/meeting';

export interface SocketResponse<S> {
  code: number;
  message: string;
  timestamp: number;
  response: S;
}

export interface GetAppIDResponse {
  app_id: string;
}

export interface GetVerifyCodePayload {
  cell_phone: string;
  country_code: string;
}

export interface VerifyLoginSms {
  cell_phone: string;
  country_code: string;
  code: string;
}

export interface VerifyLoginToken {
  meaningless?: null;
}

export interface LoginToken {
  login_token: string;
}
export interface VerifyLoginRes {
  created_at: number;
  user_id: string;
  user_name: string;
  login_token: string;
}
export interface GetAppIDPayload {
  meaningless?: null;
}

export interface JoinMeetingPayload {
  app_id: string;
  user_id: string;
  user_name: string;
  room_id: string;
  mic: boolean;
  camera: boolean;
}

export interface JoinMeetingResponse {
  token: string;
  info: MeetingInfo;
  users: MeetingUser[];
}

export interface UserPayload {
  user_id?: string;
}

export interface RecordMeetingPayload {
  users: string[];
  screen_uid: string;
}

export interface EndMeetingPayload {
  meaningless?: null;
}

export interface UserStatusChangePayload {
  user_id: string;
  status: boolean;
}

export interface HostChangePayload {
  former_host_id: string;
  host_id: string;
}

export interface HistoryVideoRecord {
  created_at: number;
  download_url: string;
  room_id: string;
  vid: string;
  video_holder: boolean;
}

export interface AudioStats {
  CodecType: string;
  End2EndDelay: number;
  MuteState: boolean;
  PacketLossRate: number;
  RecvBitrate: number;
  RecvLevel: number;
  TotalFreezeTime: number;
  TotalPlayDuration: number;
  TransportDelay: number;
}
