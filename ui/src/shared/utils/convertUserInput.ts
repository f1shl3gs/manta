import React from 'react'

export const convertUserInputValueToNumOrNaN = (value?: any): number =>
  value === '' ? NaN : Number(value)

export const convertUserInputToNumOrNaN = (
  ev: React.ChangeEvent<HTMLInputElement>
): number => convertUserInputValueToNumOrNaN(ev.target.value)
