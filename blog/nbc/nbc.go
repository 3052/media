package main

import (
   "fmt"
   "io"
   "net/http"
   "net/url"
   "os"
   "strings"
)

func main() {
   var req http.Request
   req.Header = http.Header{}
   req.Header.Add("Content-Type", "application/json")
   req.Method = "POST"
   req.ProtoMajor = 1
   req.ProtoMinor = 1
   req.URL = &url.URL{}
   req.URL.Host = "friendship.nbc.com"
   req.URL.Path = "/v3/graphql"
   req.URL.Scheme = "https"
   data = fmt.Sprintf(`
   {
     "query": %q,
     "variables": {
         "userId": "",
         "name": "saturday-night-live/video/november-15-glen-powell/9000454161",
         "app": "nbc",
         "platform": "web",
         "type": "VIDEO"
     }
   }
   `, data)
   req.Body = io.NopCloser(strings.NewReader(data))
   resp, err := http.DefaultClient.Do(&req)
   if err != nil {
      panic(err)
   }
   err = resp.Write(os.Stdout)
   if err != nil {
      panic(err)
   }
}

var data = `
fragment videoPageMetaData on VideoPageMetaData {
  mpxAccountId
  mpxGuid
  programmingType

  title
  secondaryTitle
  tertiaryTitle
  playlistTitle
  playlistMachineName
  playlistImage {
    ...image
  }
  playlistCreated
  playlistDescription
  description
  shortDescription
  gradientStart
  gradientEnd
  lightPrimaryColor
  darkPrimaryColor
  brandLightPrimaryColor
  brandDarkPrimaryColor
  seasonNumber
  episodeNumber
  airDate
  rating
  copyright
  locked
  genre
  amazonGenre
  duration
  permalink
  image
  labelBadge
  authEnds
  externalAdvertiserId
  mpxEntitlementWindows {
    availStartDateTime
    availEndDateTime
    entitlement
    device
  }
  tveEntitlementWindows {
    availStartDateTime
    availEndDateTime
    entitlement
    device
  }
  cast {
    characterFirstName
    characterLastName
    talentFirstName
    talentMiddleName
    talentLastName
  }
  v4ID
  seriesShortTitle
  seriesShortDescription
  multiPlatformLargeImage
  multiPlatformSmallImage
  urlAlias
  seriesType
  dayPart
  sunrise
  sunset
  ratingAdvisories
  width
  height
  selectedCountries
  keywords
  watchId
  numberOfEpisodes
  numberOfSeasons
  channelId
  resourceId
  mpxAdPolicy
  brandDisplayTitle
  brandMachineName
  colorBrandLogo
  whiteBrandLogo
  cuePoint
  startRecapTiming
  endRecapTiming
  startTeaserTiming
  endTeaserTiming
  startIntroTiming
  endIntroTiming
  allowSkipButtons
  skipButtonsDuration
  allowMobileWebPlayback
  canonicalUrl
  ariaLabel
  tmsId
  movieShortTitle
  headerTitle
  tuneIn
}


fragment image on Image {
  path
  width
  height
  altText
}

fragment streamPageAnalyticsAttributes on StreamPageAnalyticsAttributes {
  programmingType

  pageType
  brand {
    title
  }
  genre
  secondaryGenre
  sport
  league
  rating
  ratingAdvisories
}



fragment coreSDKPlayer on CoreSDKPlayer {
  ...component
  ...section
  CoreSDKPlayerData: data {
    ...componentData
    player {
      mpxAccountId
      mpxGuid
      programmingType
      
      v4ID
      contentType
      title
      titleId
      permalink
      description
      secondaryTitle
      secondaryGenre
      pid
      image
      airDate
      brandDisplayTitle
      brandMachineName
      colorBrandLogo
      whiteBrandLogo
      resourceId
      regionEntitlementId
      channelId
      sport
      league
      offerType
      tuneIn
      gradientStart
      gradientEnd
      tertiaryTitle
      entitlement
      locked
      endTime
      startTime
      backgroundImage
      headerTitle
      headerTitleLogo
      duration
      rating
      ratingAdvisories
      nielsenSfCode
      ariaLabel
      seasonNumber
      episodeNumber
      genre
      amazonGenre
      copyright
      dayPart
      shortDescription
      sunset
      keywords
      seriesShortTitle
      seriesUrlAlias
      lightPrimaryColor
      mpxAdPolicy
      allowMobileWebPlayback
      startRecapTiming
      endRecapTiming
      startTeaserTiming
      endTeaserTiming
      startIntroTiming
      endIntroTiming
      cuePoint
      externalAdvertiserId
      tmsId
      movieShortTitle
      allowSkipButtons
      skipButtonsDuration
      goToButtonText
      goToButtonDestination
      stationId
      callSign
      streamAccessName
      notification {
        ...notification
      }
      brandV4ID
      broadcastRegion
      contentGatingType
      availableForCredits
    }
    endCard {
      ...lazyEndCard
    }
  }
}

fragment component on Component {
  component
  meta
  treatments
}



fragment section on Section {
  logicName
  deepLinkHandle
}



fragment componentData on ComponentData {
  instanceID
}



fragment notification on Notification {
  ...component
  data {
    ...componentData
    v4ID
    machineName
    headline
    headlineColor
    message
    messageColor
    logo
    logoAltText
    portraitImage
    landscapeImage
    cta {
      ...ctaLink
    }
    dismissText
  }
  analytics {
    entityTitle
    dismissText
    placement
  }
}



fragment ctaLink on CTALink {
  ...component
  data {
    ...ctaData
  }
  analytics {
    mpxGuid
    programmingType
   
    ctaTitle
    destinationType
    destination
    brand {
      title
    }
    series
    movie
    isMovie
    videoTitle
    locked
    seasonNumber
    episodeNumber
    duration
    isPlaylist
    playlistMachineName
    playlistTitle
    isLive
    sponsorName
    isSponsoredTitle
    isSportVideo
    language
    league
    event
    sport
    pid
  }
}



fragment ctaData on CTAData {
  ...componentData
  color
  gradientStart
  gradientEnd
  text
  destinationType
  destination
  endCardMpxGuid
  endCardTagLine
  playlistMachineName
  playlistCount
  urlAlias
  isLive
  isPlaylist
  title
  secondaryTitle
  secondaryTitleTag
  isSportsVideoSlide
  image {
    ...image
  }
  pid
  contentType
  ariaLabel
}



fragment lazyEndCard on LazyEndCard {
  ...component
  ...lazyComponent
}



fragment lazyComponent on LazyComponent {
  targetComponent
  lazyData {
    ...componentData
    queryName
    queryVariables
    entryField
    path
  }
}



fragment lazyOnAirNowShelf on LazyOnAirNowShelf {
  ...component
  ...section
  ...lazyComponent
}



fragment onAirNowShelf on OnAirNowShelf {
  ...component
  ...section
  data {
    ...onAirNowList
  }
  analytics {
    itemsList
    machineName
    listTitle
  }
}



fragment onAirNowList on OnAirNowList {
  ...componentData
  machineName
  listTitle
  ariaLabel
  lastModified
  items {
    ...onAirNowTile
  }
}



fragment onAirNowTile on OnAirNowTile {
  ...component
  onAirNowTileData: data {
    ...onAirNowItem
  }
  analytics {
    isLive
    entitlement
    episodeNumber
    seasonNumber
    programTitle
    episodeTitle
    tmsId
    adobeVideoResearchTitle
    league
    isOlympics
    sport
    nielsenSfCode
    nielsenChannel
    nielsenClientId
    videoBroadcast
    brand {
      title
    }
  }
}

fragment onAirNowItem on OnAirNowItem {
  mpxGuid

  ...componentData
  v4ID
  image
  switchToNationalStream
  title
  secondaryTitle
  startTime
  endTime
  brandV4ID
  machineName
  whiteBrandLogo
  brandDisplayTitle
  brandLightPrimary
  brandDarkPrimary
  isNew
  audioDescription
  ratingWithAdvisories
  badge
  resourceId
  channelId
  nextEpisodeMpxGuid
  relativePath
  nextEpisodeRelativePath
  watchTagline
  ariaLabel
  streamAccessName
  stationId
  callSign
  contentType
  notification {
    ...notification
  }
}



fragment linksSelectableGroup on LinksSelectableGroup {
  ...component
  ...section
  data {
    ...stringSelectableComponentList
  }
  analytics {
    itemLabels
  }
}



fragment stringSelectableComponentList on StringSelectableComponentList {
  ...componentData
  initiallySelected
  itemLabelsConfig {
    ...itemLabelsConfigItem
  }
  itemLabelsTitle
  optionalTitle: title
  gradientStart
  items {
    ...shelf
    ...lazyShelf
    ...shelfGroup
    ...lazyShelfGroup
    ...grid
    ...stack
    ...lazyStack
    ...placeholder
    ...nestedLinksSelectableGroup
    ...navigationMenuGroup
    ...navigationItem
  }
}



fragment itemLabelsConfigItem on ItemLabelsConfigItem {
  itemLabel
  menuItemType
  index
  isSelected
}



fragment shelf on Shelf {
  ...component
  ...section
  data {
    ...tileList
  }
  analytics {
    isPlaylist
    playlistMachineName
    listTitle
    isSponsoredContent
    sponsorName
    isMixedTiles
    machineName
    itemsList
  }
}



fragment tileList on TileList {
  ...componentData
  machineName
  playlistMachineName
  listTitle
  ariaLabel
  listTitleImage
  sponsorLogo
  sponsorName
  sponsorLogoAltText
  sponsorFreeWheelTrackingURL
  sponsorThirdPartyTrackingURL
  lastModified
  items {
    ...videoTile
    ...seriesTile
    ...movieTile
    ...brandTile
    ...personTile
    ...featureTile
    ...playlistTile
    ...marketingBand
    ...upcomingLiveSlideTile
    ...onAirNowTile
    ...genreTile
    ...upcomingLiveTile
    ...replayTile
  }
  moreItems {
    ...lazyShelf
    ...lazyGrid
    ...lazyStack
  }
  viewAllCta {
    ...ctaLink
  }
}



fragment videoTile on VideoTile {
  ...component
  data {
    ...videoItem
  }
  analytics {
    mpxGuid
    programmingType
   
    brand {
      title
    }
    series
    title
    episodeNumber
    seasonNumber
    locked
    duration
    movie
    genre
    sport
    league
    language
    event
    permalink
    game
    gameType
    gamesList
    isOlympics
  }
}

fragment videoItem on VideoItem {
  mpxAccountId
  mpxGuid
  programmingType

  ...componentData
  ...item
  secondaryTitleTag
  locked
  episodeNumber
  seasonNumber
  airDate
  percentViewed
  permalink
  lastWatched
  duration
  genre
  rating
  lightPrimaryColor
  darkPrimaryColor
  seriesShortTitle
  movieShortTitle
  whiteBrandLogo
  colorBrandLogo
  brandDisplayTitle
  mpxAdPolicy
  resourceId
  channelId
  rating
  externalAdvertiserId
  ariaLabel
  longDescription
  ctaText
  ctaTextColor
  brandMachineName
  durationBadge
  gradientStart
  gradientEnd
  isOlympics
}



fragment item on Item {
  v4ID
  title
  secondaryTitle
  tertiaryTitle
  description
  image
  gradientStart
  gradientEnd
  labelBadge
  lastModified
}



fragment seriesTile on SeriesTile {
  ...component
  data {
    ...seriesItem
  }
  analytics {
    series
    brand {
      title
    }
    genre
    sport
    league
  }
}



fragment seriesItem on SeriesItem {
  ...componentData
  ...item
  urlAlias
  posterImage
  whiteBrandLogo
  brandDisplayTitle
  ariaLabel
  gradientStart
  gradientEnd
}



fragment movieTile on MovieTile {
  ...component
  data {
    ...movieItem
  }
  analytics {
    isTrending
    movie
    brand {
      title
    }
    genre
  }
}



fragment movieItem on MovieItem {
  ...componentData
  ...item
  urlAlias
  posterImage
  image
  lightPrimaryColor
  darkPrimaryColor
  whiteBrandLogo
  colorBrandLogo
  brandDisplayTitle
  landscapePreview
  portraitPreview
  rating
  ariaLabel
}



fragment brandTile on BrandTile {
  ...component
  data {
    ...brandItem
  }
  analytics {
    brand {
      title
    }
  }
}



fragment brandItem on BrandItem {
  ...componentData
  v4ID
  displayTitle
  machineName
  lightPrimaryColor
  darkPrimaryColor
  colorBrandLogo
  whiteBrandLogo
  horizontalPreview
  staticPreviewImage
  ariaLabel
  routeToLiveStream
}



fragment personTile on PersonTile {
  ...component
  data {
    ...personItem
  }
}



fragment personItem on PersonItem {
  ...componentData
  title
  secondaryTitle
  personImage: image {
    path
  }
  machineName
  roleName
  roleMachineName
  ariaLabel
  titleUrlAlias
  gradientStart
  seasonNumber
}



fragment featureTile on FeatureTile {
  ...component
  data {
    ...featureItem
  }
  analytics {
    series
    brand {
      title
    }
    playlistMachineName
    listTitle
  }
}



fragment featureItem on FeatureItem {
  ...componentData
  ...item
  seriesShortTitle
  image
  brandDisplayTitle
  whiteBrandLogo
  colorBrandLogo
  destinationType
  destination
  playlistMachineName
  ariaLabel
}



fragment playlistTile on PlaylistTile {
  ...component
  data {
    ...playlistItem
  }
  analytics {
    brand {
      title
    }
    playlistMachineName
    listTitle
  }
}



fragment playlistItem on PlaylistItem {
  ...componentData
  ...item
  brandDisplayTitle
  whiteBrandLogo
  colorBrandLogo
  destination
  destType: destinationType
  playlistMachineName
  externalAdvertiserId
  playlistCount
  ariaLabel
}



fragment marketingBand on MarketingBand {
  ...component
  data {
    ...marketingBandData
  }
  analytics {
    series
    brand {
      title
    }
  }
}



fragment marketingBandData on MarketingBandData {
  ...componentData
  primaryImage
  compactImage
  link
  seriesShortTitle
  ariaLabel
}



fragment upcomingLiveSlideTile on UpcomingLiveSlideTile {
  ...component
  data {
    ...upcomingLiveSlideTileData
  }
  analytics {
    analyticsType
    ctaLiveTitle
    ctaUpcomingTitle
    ctaNotInPackageTitle
    isLiveCallout
    programType
    genre
    secondaryGenre
    listOfGenres
    title
    secondaryTitle
    league
    sport
    videoBroadcast
    nielsenClientId
    nielsenChannel
    nielsenSfCode
    isOlympics
    adobeVideoResearchTitle
    brand {
      title
    }
    seasonNumber
    permalink
    videoTitle
  }
}



fragment upcomingLiveSlideTileData on UpcomingLiveSlideTileData {
  ...componentData
  v4ID
  title
  secondaryTitle
  description
  image
  gradientStart
  gradientEnd
  liveBadge
  upcomingBadge
  titleColor
  titleLogo
  secondaryTitleColor
  descriptionColor
  compactImage
  landscapePreview
  brandDisplayTitle
  whiteBrandLogo
  colorBrandLogo
  customerPlayableDate
  startTime
  endTime
  liveAriaLabel
  upcomingAriaLabel
  liveCtaText
  upcomingCtaText
  notInPackageCtaText
  resourceId
  channelId
  machineName
  streamAccessName
  directToLiveThreshold
  stationId
  contentType
  pid
  relativePath
  upcomingModal {
    ...upcomingModal
  }
  notification {
    ...notification
  }
  programType
  genre
  secondaryGenre
  sport
  league
  callSign
}



fragment upcomingModal on UpcomingModal {
  ...component
  data {
    ...upcomingModalData
  }
  analytics {
    modalName
    modalType
    dismissText
  }
}



fragment upcomingModalData on UpcomingModalData {
  machineName
  title
  description
  ctaText
  dismissText
  lastMinuteModalLifespan
  countdownDayLabel
  countdownHourLabel
  countdownMinLabel
  customerPlayableDate
  startTime
  backgroundImage
  contentType
  pid
  streamAccessName
}



fragment genreTile on GenreTile {
  ...component
  data {
    title
    image
    gradientStart
    gradientEnd
    ctaLink {
      ...ctaLink
    }
    ariaLabel
  }
  analytics {
    itemClickedName
    itemClickedType
    machineName
  }
}



fragment upcomingLiveTile on UpcomingLiveTile {
  ...component
  data {
    ...upcomingLiveItem
  }
  analytics {
    secondaryGenre
    listOfGenres
    entitlement
    locked
    league
    sport
    brand {
      title
    }
    seasonNumber
    permalink
    pid
    secondaryTitle
  }
}



fragment upcomingLiveItem on UpcomingLiveItem {
  relativePath
  title
  secondaryTitle
  liveBadge
  upcomingBadge
  image
  startTime
  endTime
  whiteBrandLogo
  brandDisplayTitle
  liveAriaLabel
  upcomingAriaLabel
  upcomingModal {
    ...upcomingModal
  }
  streamAccessName
  directToLiveThreshold
  contentType
  pid
  notification {
    ...notification
  }
  sport
  league
  audioLanguage
  tertiaryTitle
  isMedalSession
  isOlympics
}



fragment replayTile on ReplayTile {
  ...component
  replayTileData: data {
    ...replayTileData
  }
  analytics {
    programmingType
   
    analyticsType
    title
    brand {
      title
    }
    genre
    nielsenSfCode
    sport
    league
    event
    secondaryGenre
    listOfGenres
    entitlement
    locked
    duration
    pid
    game
    gameType
    gamesList
    isOlympics
  }
}

fragment replayTileData on ReplayTileData {
  programmingType

  ...componentData
  v4ID
  ariaLabel
  brandDisplayTitle
  colorBrandLogo
  image
  pid
  relativePath
  secondaryTitle
  title
  whiteBrandLogo
  tertiaryTitle
  labelBadge
  locked
  listTitle
  isOlympics
  audioLanguage
}



fragment lazyShelf on LazyShelf {
  ...component
  ...section
  ...lazyComponent
}



fragment lazyGrid on LazyGrid {
  ...component
  ...section
  ...lazyComponent
}



fragment lazyStack on LazyStack {
  ...component
  ...section
  ...lazyComponent
}



fragment shelfGroup on ShelfGroup {
  ...component
  ...section
  data {
    ...shelfList
  }
}



fragment shelfList on ShelfList {
  ...componentData
  listTitle
  items {
    ...shelf
  }
}



fragment lazyShelfGroup on LazyShelfGroup {
  ...component
  ...lazyComponent
  ...section
}



fragment grid on Grid {
  ...component
  ...section
  data {
    ...tileList
  }
  analytics {
    listTitle
    playlistMachineName
    isSponsoredContent
    sponsorName
  }
}



fragment stack on Stack {
  ...component
  ...section
  data {
    ...tileList
  }
  analytics {
    playlistMachineName
    listTitle
    isSponsoredContent
    sponsorName
  }
}



fragment placeholder on Placeholder {
  ...component
  deepLinkHandle
  data {
    ...componentData
    machineName
    placeholderType
    queryVariables
    queryName
    entryField
    path
  }
}



fragment nestedLinksSelectableGroup on NestedLinksSelectableGroup {
  ...component
  deepLinkHandle
  data {
    ...nestedStringSelectableComponentList
  }
  analytics {
    itemLabels
  }
}



fragment nestedStringSelectableComponentList on NestedStringSelectableComponentList {
  ...componentData
  initiallySelected
  itemLabelsConfig {
    ...itemLabelsConfigItem
  }
  itemLabelsTitle
  optionalTitle: title
  gradientStart
  items {
    ...shelf
    ...lazyShelf
    ...shelfGroup
    ...lazyShelfGroup
    ...grid
    ...stack
    ...lazyStack
    ...placeholder
  }
}



fragment navigationMenuGroup on NavigationMenuGroup {
  ...component
  data {
    ...navigationMenuList
  }
  deepLinkHandle
  analytics {
    analyticsType
  }
}



fragment navigationMenuList on NavigationMenuList {
  ...componentData
  listTitle
  listTitleImage
  items {
    ...navigationMenu
  }
}



fragment navigationMenu on NavigationMenu {
  ...component
  ...section
  data {
    ...componentData
    ...navigationMenuData
  }
}



fragment navigationMenuData on NavigationMenuData {
  ...componentData
  titleMenu: title
  items {
    ...navigationItem
  }
}



fragment navigationItem on NavigationItem {
    data {
    ...navigationItemData
  }
  analytics {
    analyticsType
    itemClickedName
  }
}



fragment navigationItemData on NavigationItemData {
  title
  defaultLogo
  ariaLabel
  destination
  isLive
}



fragment expiredVideo on VideoDetailsExpired {
  ...component
  ...section
  data {
    ...componentData
    videoMeta {
      title
      secondaryTitle
      description
      image
    }
  }
}



fragment streamNotFound on StreamNotFound {
  ...component
  ...section
  data {
    ...componentData
    backgroundImage
  }
}

query page(
  $id: ID
  $name: String!
  $queryName: QueryNames
  $type: PageType!
  $subType: PageSubType
  $nationalBroadcastType: String
  $userId: String!
  $platform: SupportedPlatforms!
  $device: String
  $profile: JSON
  $timeZone: String
  $deepLinkHandle: String
  $app: NBCUBrands!
  $nbcAffiliateName: String
  $telemundoAffiliateName: String
  $language: Languages
  $playlistMachineName: String
  $mpxGuid: String
  $authorized: Boolean
  $minimumTiles: Int
  $endCardMpxGuid: String
  $endCardTagLine: String
  $seasonNumber: Int
  $creditMachineName: String
  $roleMachineName: String
  $originatingTitle: String
  $isDayZero: Boolean
) {
  page(
    id: $id
    name: $name
    type: $type
    subType: $subType
    nationalBroadcastType: $nationalBroadcastType
    userId: $userId
    queryName: $queryName
    platform: $platform
    device: $device
    profile: $profile
    timeZone: $timeZone
    deepLinkHandle: $deepLinkHandle
    app: $app
    nbcAffiliateName: $nbcAffiliateName
    telemundoAffiliateName: $telemundoAffiliateName
    language: $language
    playlistMachineName: $playlistMachineName
    mpxGuid: $mpxGuid
    authorized: $authorized
    minimumTiles: $minimumTiles
    endCardMpxGuid: $endCardMpxGuid
    endCardTagLine: $endCardTagLine
    seasonNumber: $seasonNumber
    creditMachineName: $creditMachineName
    roleMachineName: $roleMachineName
    originatingTitle: $originatingTitle
    isDayZero: $isDayZero
  ) {
    id
    pageType
    name
    metadata {
      __typename
      ...videoPageMetaData
    }
    analytics {
      ...streamPageAnalyticsAttributes
    }
    data {
      sections {
        ...coreSDKPlayer
        ...lazyOnAirNowShelf
        ...onAirNowShelf
        ...linksSelectableGroup
        ...shelf
        ...stack
        ...expiredVideo
        ...streamNotFound
      }
      menu {
        ...navigationMenu
      }
    }
  }
}
`
