import { util } from '@aws-appsync/utils';

export function request(ctx) {
  return {
    payload: ctx.args
  };
}

export function request(ctx) {
  return ctx.result;
}
