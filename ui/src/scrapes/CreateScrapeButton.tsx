import React, {FunctionComponent} from 'react'
import {Button, ComponentColor, IconFont} from '@influxdata/clockface'
import {useDispatch} from 'react-redux'
import {createScrape} from 'src/scrapes/actions/thunk'

const CreateScrapeButton: FunctionComponent = () => {
  const dispatch = useDispatch()
  const handleClick = () => {
    dispatch(createScrape())
  }

  return (
    <Button
      text={'Create'}
      testID={'create-scrape--button'}
      color={ComponentColor.Primary}
      icon={IconFont.Plus_New}
      onClick={handleClick}
    />
  )
}

export default CreateScrapeButton
