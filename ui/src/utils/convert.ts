export const convertUserInputToNumOrNaN = (
  e: React.ChangeEvent<HTMLInputElement>
) => (e.target.value === '' ? NaN : Number(e.target.value))
