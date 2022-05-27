import React from "react";

function TradeCard() {
  return (
    <div className="mx-auto flex w-full max-w-screen-2xl flex-col justify-center bg-white rounded-2xl shadow-xl shadow-slate-300/60">
      <div className="p-6">
        <p>
          <small className="text-s mr-1 pointer-events-none">Ticker:</small>
          <a className="font-bold text-indigo-500 hover:text-indigo-900" href="/">MDRR</a>
        </p>
        <h1 className="text-2xl font-medium text-slate-600 pb-2">Company</h1>
        <p className="text-sm tracking-tight font-light text-slate-400 leading-6">
          Dodge is an American brand of automobiles and a division of Stellantis, based in Auburn Hills, Michigan..
        </p>
      </div>
    </div>
  );
}

export default TradeCard;
