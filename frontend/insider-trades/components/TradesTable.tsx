import React, { useEffect, useRef, useState } from "react";
import { useAppDispatch, useAppSelector } from "store/hooks";
import { fetchTradeViews } from "store/tradeViewsSlice";
import { TradeView } from "types/tradeView";

interface Options {
  callback: () => Promise<unknown>;
  element: HTMLElement | null;
}

const useInfiniteScroll = ({ callback, element }: Options) => {
  const [isFetching, setIsFetching] = useState(false);
  const observer = useRef<IntersectionObserver>();

  useEffect(() => {
    if (!element) {
      return;
    }

    observer.current = new IntersectionObserver(
      (entries) => {
        if (!isFetching && entries[0].isIntersecting) {
          setIsFetching(true);
          callback().finally(() => setIsFetching(false));
        }
      },
      {
        rootMargin: "-100px",
      }
    );
    observer.current.observe(element);

    return () => observer.current?.disconnect();
  }, [callback, isFetching, element]);

  return isFetching;
};

export const TradesTable = () => {
  const dispatch = useAppDispatch();
  const tradeViews = useAppSelector((state) => state.tradeViews);

  useEffect(() => {
    dispatch(fetchTradeViews({ refresh: true }));
  }, [dispatch]);

  const getTradeTypeClass = (transactionTypeName: string) => {
    const buyType = "BUY";
    const sellTypeClass = "trade-type type-sell";
    const buyTypeClass = "trade-type type-buy";

    if (transactionTypeName === buyType) {
      return buyTypeClass;
    } else {
      return sellTypeClass;
    }
  };
  const USDFormatter = (minFraction: number, maxFraction: number) => {
    return new Intl.NumberFormat("en-US", {
      style: "currency",
      currency: "USD",

      minimumFractionDigits: minFraction,
      maximumFractionDigits: maxFraction,
    });
  };

  return (
    <div>
      <table className="trades-table rounded-lg">
        <thead>
          <tr>
            <th>Ticker</th>
            <th>Company</th>
            <th>Insider</th>
            <th>Type</th>
            <th>Total shares</th>
            <th>Average price</th>
            <th>Total value</th>
            <th>Reported on</th>
          </tr>
        </thead>
        <tbody>
          {tradeViews.tradeViews.map((trade: TradeView) => {
            if (trade == undefined) {
              return;
            }
            return (
              <tr key={trade.ID}>
                <td className="ticker">{trade.CompanyTicker}</td>
                <td>{trade.CompanyName}</td>
                <td>{trade.InsiderName}</td>
                <td>
                  <p className={getTradeTypeClass(trade.TransactionTypeName)}>
                    {trade.TransactionTypeName}
                  </p>
                </td>
                <td>{trade.TotalShares}</td>
                <td>{USDFormatter(3, 3).format(trade.AveragePrice)}</td>
                <td>{USDFormatter(0, 0).format(trade.TotalValue)}</td>
                <td>{trade.ReportedOn}</td>
              </tr>
            );
          })}
        </tbody>
      </table>
    </div>
  );
};
