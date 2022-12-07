import {
  AutoInput,
  AutoInputMode,
  Form,
  Input,
  InputType,
} from '@influxdata/clockface'
import React, {FunctionComponent, useState} from 'react'
import {MAX_DECIMAL_PLACES, MIN_DECIMAL_PLACES} from 'src/shared/constants/dashboard'

// Utils
import {convertUserInputToNumOrNaN} from 'src/shared/utils/convertUserInput'

interface Props {
  isEnforced: boolean
  digits: number
  update: (obj: any) => void
}

const DecimalPlaces: FunctionComponent<Props> = ({
  isEnforced,
  digits,
  update,
}) => {
  const [decimalPlaces, setDecimalPlaces] = useState(digits)

  const setDigits = (updated: number | null) => {
    setDecimalPlaces(updated)

    if (!Number.isNaN(updated)) {
      update({
        decimalPlaces: {
          isEnforced,
          digits: updated,
        },
      })
    }
  }

  const handleChangeMode = (mode: AutoInputMode): void => {
    if (mode === AutoInputMode.Auto) {
      setDigits(null)
    } else {
      setDigits(2)
    }
  }

  return (
    <Form.Element label={'Decimal Places'}>
      <AutoInput
        inputComponent={
          <Input
            name={'decimal-places'}
            placeholder={'Enter a number'}
            onChange={ev => setDigits(convertUserInputToNumOrNaN(ev))}
            value={decimalPlaces}
            min={MIN_DECIMAL_PLACES}
            max={MAX_DECIMAL_PLACES}
            type={InputType.Number}
          />
        }
        onChangeMode={handleChangeMode}
        mode={
          typeof decimalPlaces === 'number'
            ? AutoInputMode.Custom
            : AutoInputMode.Auto
        }
      />
    </Form.Element>
  )
}

export default DecimalPlaces
