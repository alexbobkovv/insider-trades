import { createSlice, PayloadAction, createAsyncThunk } from "@reduxjs/toolkit";
import { TradeView } from "types/tradeView";

export type TradeViewsState = {
  tradeViews: TradeView[];
  nextCursor: string | null;
  status: string | null;
  isError: boolean;
  errorMessage: string | null;
};

export type TradeViewsResponse = {
  tradeViews: TradeView[];
};

const initialState: TradeViewsState = {
  tradeViews: [],
  nextCursor: null,
  status: null,
  isError: false,
  errorMessage: null,
};

export type fetchOptions = {
  refresh?: Boolean
  nextCursor?: string
}

type requestHeaders = {
  accept: string
  cursor?: string
}

export const fetchTradeViews = createAsyncThunk<TradeViewsResponse, fetchOptions, {rejectValue: string}>(
  "tradeViews/fetchTradeViews",
  async (options: fetchOptions, {rejectWithValue}) => {
    const defaultTradeViewsURL = "http://localhost:8082/api-gateway/v1/trade-views"
    var tradeViewsURL = process.env.TRADES_GATEWAY_URL
    if (tradeViewsURL == undefined) {
      tradeViewsURL = defaultTradeViewsURL
    }

    try {
      var reqHeaders: requestHeaders = {
        accept: "application/json",
      }
      if (options.nextCursor) {
        reqHeaders.cursor = options.nextCursor;
      }

      const response = await fetch(tradeViewsURL, {
        method: "GET",
        headers: reqHeaders,
      });

      if (!response.ok) {
        const message: string = `An error has occured: ${response.status}`;
        return rejectWithValue(message)
      }

      const trades: TradeView[] = await response.json();

      return {tradeViews: trades} as TradeViewsResponse;

    } catch (err) {
      const internalServerError = 500;
      const message: string = `An error has occured: ${internalServerError}`;

      return rejectWithValue(message)
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
  },
  extraReducers: (builder) => {
    builder.addCase(fetchTradeViews.fulfilled, (state, action) => {
      const tradesResponse = action.payload?.tradeViews;

      if (tradesResponse == undefined) {
        return
      }
      if (tradesResponse && tradesResponse.length > 0) {
        if (action.meta.arg.refresh) {
          state.tradeViews = tradesResponse
        } else {
          state.tradeViews.push(...tradesResponse);
        }
        state.status = "Resolved"
        state.isError = false;
      }
    }),

    builder.addCase(fetchTradeViews.pending, (state) => {
      state.status = "Loading.."
      state.isError = false;
      state.errorMessage = null;
    }),

    builder.addCase(fetchTradeViews.rejected, (state, action) => {
      state.status = "Rejected";
      state.isError = true;
      if (action.payload != undefined) {
        state.errorMessage = action.payload;
      } else {
        state.errorMessage = "unknown error"
      }
    })
  },
});

export const { appendTradeViews } = tradeViewsSlice.actions;

export default tradeViewsSlice.reducer;
