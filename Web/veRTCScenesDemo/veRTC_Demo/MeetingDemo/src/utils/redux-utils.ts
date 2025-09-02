import { Action } from 'dva-model-creator';
import { ImmerReducer, Dispatch } from '@@/plugin-dva/connect';

export const setFieldsReducer = <T extends Record<string, any>>(): ImmerReducer<
  T,
  Action<Partial<T>>
> => {
  return (state: T, action: Action<Partial<T>>) => {
    const { payload } = action;
    for (const key in payload) {
      if (payload[key]) {
        state[key] = payload[key] as T[Extract<keyof T, string>];
      }
    }
  };
};

export const setFieldReducer = <T extends any>(
  state: T,
  key: keyof T
): ImmerReducer<T, Action<T[typeof key]>> => {
  return (state: T, action: Action<T[typeof key]>) => {
    const { payload } = action;
    state[key] = payload;
  };
};

/**
 * dispatch with a promise returned
 *
 * @param {Dispatch<P>} originDispatch
 * @return {Dispatch<P>}
 */
export const asyncDispatch = <P extends any>(
  originDispatch: Dispatch<P>
): Dispatch<P> => {
  return (action) => {
    return new Promise((res, rej) => {
      originDispatch({ ...action, resolve: res, reject: rej });
    });
  };
};
