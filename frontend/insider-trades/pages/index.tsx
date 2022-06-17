import { TradesTable } from "components/TradesTable";
import type { NextPage } from "next";
import { useAppSelector } from "store/hooks";

const Home: NextPage = () => {
  const tradeViews = useAppSelector((state) => state.tradeViews);

  return (
    <main className="w-full text-center container">
      {!tradeViews.isError ? (
        <section className="trades">
          <h1 className="mb-10">Recent insider trades</h1>
          <TradesTable />
        </section>
      ) : (
        <h2 className="text-2xl">Failed to load trades</h2>
      )}
    </main>
  );
};

export default Home;
