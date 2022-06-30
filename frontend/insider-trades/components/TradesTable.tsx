import React, { useEffect, useRef, useState } from "react";
import { useAppDispatch, useAppSelector } from "store/hooks";
import { fetchTradeViews } from "store/tradeViewsSlice";
import { TradeView } from "types/tradeView";
import { ShowMoreButton } from "./buttons/ShowMoreButton";

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
  const lastTradeViewElementRef = useRef<HTMLTableRowElement | null>(null);

  const appendTradeViews = (limit: number, refresh: boolean) => {
    dispatch(fetchTradeViews({ nextCursor: tradeViews.nextCursor, refresh: refresh, limit: limit}))
  }

  if (lastTradeViewElementRef) {
    const node = lastTradeViewElementRef.current

    useInfiniteScroll({callback: () => dispatch(fetchTradeViews({ nextCursor: tradeViews.nextCursor, refresh: false, limit: 20})), element: node as HTMLElement})
  }

  useEffect(() => {
    dispatch(fetchTradeViews({ refresh: true, limit: 20}));
  }, []);

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

  const handleShowMoreTrades = () => {
    appendTradeViews(20, false)
  }

  return (
    <div className="flex flex-col">
      <table className="trades-table mb-8 rounded-lg">
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
          {tradeViews.tradeViews.map((trade: TradeView, i: number) => {
            if (trade == undefined) {
              return;
            }
            var trRef = null
            if (i == tradeViews.tradeViews.length - 1) {
              trRef = lastTradeViewElementRef 
            }

            return (
              <tr ref={trRef} key={trade.ID}>
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
      {(!tradeViews.isLastPage && <ShowMoreButton text="Show more" onClickHandler={handleShowMoreTrades}/>)}
    </div>
  );
};
