import { configureStore } from "@reduxjs/toolkit";
import { enableMapSet } from "immer";
import tradeViewsReducer from "./tradeViewsSlice";

enableMapSet();

const store = configureStore({
  reducer: {
    tradeViews: tradeViewsReducer,
  },
});

export default store;

export type RootState = ReturnType<typeof store.getState>;
export type AppDispatch = typeof store.dispatch;
