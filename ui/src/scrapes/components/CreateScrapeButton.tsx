// Libraries
import React, {FunctionComponent} from 'react'
import {useDispatch} from 'react-redux'

// Component
import {Button, ComponentColor, IconFont} from '@influxdata/clockface'

// Actions
import {createScrape} from 'src/scrapes/actions/thunks'

const CreateScrapeButton: FunctionComponent = () => {
  const dispatch = useDispatch()
  const handleClick = () => {
    dispatch(createScrape())
  }

  return (
    <Button
      text={'Create Scrape'}
      testID={'create-scrape--button'}
      color={ComponentColor.Primary}
      icon={IconFont.Plus_New}
      onClick={handleClick}
    />
  )
}

export default CreateScrapeButton
