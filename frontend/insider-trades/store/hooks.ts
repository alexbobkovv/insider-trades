import { useDispatch, useSelector, TypedUseSelectorHook } from "react-redux";
import { AppDispatch, RootState } from "./store";

// export const useAppSelector: createSelectorHook<RootState> = useSelector;

export const useAppDispatch = () => useDispatch<AppDispatch>()
export const useAppSelector: TypedUseSelectorHook<RootState> = useSelector
