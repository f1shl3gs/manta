// Libraries
import React, {FunctionComponent} from 'react'
import {useDispatch} from 'react-redux'

// Components
import {Button, IconFont, ComponentColor} from '@influxdata/clockface'

// Actions
import {push} from '@lagunovsky/redux-react-router'

const CreateSecretButton: FunctionComponent = () => {
  const dispatch = useDispatch()
  const handleCraete = () => {
    dispatch(push(`${window.location.pathname}/new`))
  }

  return (
    <Button
      text={'Create Secret'}
      icon={IconFont.Plus_New}
      color={ComponentColor.Primary}
      testID={'create-secret--button'}
      onClick={handleCraete}
    />
  )
}

export default CreateSecretButton
