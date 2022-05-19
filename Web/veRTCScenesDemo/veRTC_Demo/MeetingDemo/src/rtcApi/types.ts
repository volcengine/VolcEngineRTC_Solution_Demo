export type ICreateStreamRes = {
  code: number;
  msg: string;
  devicesStatus: {
    video: number;
    audio: number;
  };
}
