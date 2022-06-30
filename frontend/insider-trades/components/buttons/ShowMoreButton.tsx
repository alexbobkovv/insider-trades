import React, { MouseEventHandler } from 'react'

interface Props {
  text: string
  onClickHandler: MouseEventHandler
}

export const ShowMoreButton = ({ text, onClickHandler }: Props) => {
  return (
    <button className="py-2 mb-2 border-2 border-slate-200 rounded-md hover:border-slate-400 active:bg-slate-200" onClick={onClickHandler}>{text}</button>
  )
}
