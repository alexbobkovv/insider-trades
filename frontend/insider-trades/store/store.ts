import { configureStore } from "@reduxjs/toolkit";
import tradeViewsReducer from "./tradeViewsSlice";

const store = configureStore({
  reducer: {
    tradeViews: tradeViewsReducer,
  },
});

export default store;

export type RootState = ReturnType<typeof store.getState>;
export type AppDispatch = typeof store.dispatch;
