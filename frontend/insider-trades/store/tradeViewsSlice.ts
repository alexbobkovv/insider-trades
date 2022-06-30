import { createSlice, PayloadAction, createAsyncThunk } from "@reduxjs/toolkit";
import { TradeView } from "types/tradeView";
import { RootState } from "./store";

// TODO Write custom error types and handlers
export type TradeViewsState = {
  tradeViews: TradeView[];
  nextCursor: string | null;
  prevCursor: string | null;
  status: string | null;
  isBusy: boolean;
  isError: boolean;
  isLastPage: boolean;
  errorMessage: string | null;
};

export type TradeViewsResponse = {
  tradeViews: TradeView[] | null;
  nextCursor: string | null;
};

// TODO Add status and error messages
const initialState: TradeViewsState = {
  tradeViews: [],
  nextCursor: null,
  prevCursor: null,
  status: null,
  isBusy: false,
  isError: false,
  isLastPage: false,
  errorMessage: null,
};

export type fetchOptions = {
  refresh?: boolean;
  nextCursor?: string | null;
  limit?: number;
};

function sleep(ms: number) {
  return new Promise((resolve) => setTimeout(resolve, ms));
}

// TODO Fix last cursor on backend, it should't point to empty page or it should return some json data that indicates the last page
export const fetchTradeViews = createAsyncThunk<
  TradeViewsResponse,
  fetchOptions,
  { rejectValue: string; state: RootState }
>(
  "tradeViews/fetchTradeViews",
  async (options, { rejectWithValue, getState, dispatch }) => {
    var state = getState().tradeViews;
    if (state.isLastPage) {
      return { tradeViews: null, nextCursor: null } as TradeViewsResponse;
    }

    if (state.isBusy) {
      let attempts = 100;
      while (attempts > 0) {
        const state = getState().tradeViews;
        await sleep(50);
        if (!state.isBusy) {
          break;
        }
        attempts--;
      }

      if (attempts <= 0) {
        return { tradeViews: null, nextCursor: null } as TradeViewsResponse;
      }
    }

    dispatch(lockTradeViews());
    var state = getState().tradeViews;

    if (!options.refresh && !state.nextCursor) {
      return { tradeViews: null, nextCursor: null } as TradeViewsResponse;
    }

    if (state.nextCursor && state.nextCursor == state.prevCursor) {
      return { tradeViews: null, nextCursor: null } as TradeViewsResponse;
    }

    const defaultTradeViewsURL = new URL(
      "http://localhost:8082/api-gateway/v1/trade-views"
    );

    let tradeViewsURL: URL;
    const url = process.env.TRADES_GATEWAY_URL;
    if (url == undefined) {
      tradeViewsURL = defaultTradeViewsURL;
    } else {
      tradeViewsURL = new URL(url);
    }

    if (state.nextCursor) {
      tradeViewsURL.searchParams.set("cursor", state.nextCursor);
    }

    if (options.limit) {
      tradeViewsURL.searchParams.set("limit", options.limit.toString());
    }

    try {
      const reqHeaders = {
        Accept: "application/json",
      };

      const response = await fetch(tradeViewsURL.toString(), {
        method: "GET",
        headers: reqHeaders,
      });

      if (!response.ok) {
        const message = `An error has occured: ${response.status}`;
        return rejectWithValue(message);
      }

      const trades: TradeView[] = await response.json();

      return {
        tradeViews: trades,
        nextCursor: response.headers.get("X-next-cursor"),
      } as TradeViewsResponse;
    } catch (err) {
      const internalServerError = 500;
      const message = `An error has occured: ${internalServerError}`;

      return rejectWithValue(message);
    }
  }
);

export const tradeViewsSlice = createSlice({
  name: "tradeViews",
  initialState,
  reducers: {
    appendTradeViews: (
      state: TradeViewsState,
      action: PayloadAction<TradeViewsResponse>
    ) => {
      const tradesResponse = action.payload.tradeViews;

      if (tradesResponse && tradesResponse.length > 0) {
        state.tradeViews.push(...tradesResponse);
      }
    },
    lockTradeViews: (state: TradeViewsState) => {
      state.isBusy = true;
    },
  },
  extraReducers: (builder) => {
    builder.addCase(fetchTradeViews.fulfilled, (state, action) => {
      const tradeViewsResponse = action.payload?.tradeViews;

      if (tradeViewsResponse == undefined || tradeViewsResponse == null) {
        state.isBusy = false;
        return;
      }

      if (tradeViewsResponse.length === 0) {
        state.isBusy = false;
        state.isLastPage = true;
        return;
      }

      if (action.meta.arg.refresh) {
        state.nextCursor = action.payload.nextCursor;
        state.tradeViews = tradeViewsResponse;
      } else {
        state.prevCursor = state.nextCursor;
        state.nextCursor = action.payload.nextCursor;
        state.tradeViews.push(...tradeViewsResponse);
      }
      state.status = "Resolved";
      state.isError = false;
      state.isBusy = false;
    }),
      builder.addCase(fetchTradeViews.pending, (state) => {
        state.status = "Loading..";
        state.isError = false;
        state.errorMessage = null;
      }),
      builder.addCase(fetchTradeViews.rejected, (state, action) => {
        state.status = "Rejected";
        if (state.tradeViews.length < 1) {
          state.isError = true;
        }

        if (action.payload != undefined) {
          state.errorMessage = action.payload;
        } else {
          state.errorMessage = "unknown error";
        }
      });
  },
});

export const { appendTradeViews, lockTradeViews } = tradeViewsSlice.actions;

export default tradeViewsSlice.reducer;
