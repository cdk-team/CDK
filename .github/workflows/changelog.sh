set -x
set +e


DATE_STRING=`date -u +"%Y-%m-%d"`
SHA256_TEXT_BODY=`cd bin/ && shasum -a 256 * | tr -s '  ' '|'`

TAG_VERSION=`echo "$GITHUB_REF" | sed -e 's/refs\/tags\///' || "HEAD"`
LAST_TAG="$TAG_VERSION"
PREVIOUS_TAG=`git for-each-ref --sort='-authordate' --format="%(refname:short)" | grep -E "^v[0-9]+\.[0-9]+\.[0-9]+$" | sed -n 2p`

# generate base log of exploit evaluate and tool
exploit=`git log "${LAST_TAG}...${PREVIOUS_TAG}" --pretty=format:%s -- "pkg/exploit/" | grep -viE ^merge` || :
evaluate=`git log "${LAST_TAG}...${PREVIOUS_TAG}" --pretty=format:%s -- "pkg/evaluate/" | grep -viE ^merge` || :
tool=`git log "${LAST_TAG}...${PREVIOUS_TAG}" --pretty=format:%s -- "pkg/tool/" | grep -viE ^merge` || :

# log infomation calculate
add_before=`echo "$exploit\n$evaluate\n$tool" | uniq`
all_commit_message=`git log "${LAST_TAG}...${PREVIOUS_TAG}" --pretty=format:%s | grep -viE ^merge` || :
other=`diff -u <(echo "$add_before") <(echo "$all_commit_message") | grep -E "^\+[a-zA-Z]" | cut -c 2-` || :

# add changelog format
[[ $exploit = *[^[:space:]]* ]] && exploit=`echo "$exploit" | awk '{print toupper(substr($0,1,1))""substr($0,2)}' | sed -e 's/^/* /'` && exploit=$'### :bomb: Exploits \n\n'"$exploit"
[[ $evaluate = *[^[:space:]]* ]] && evaluate=`echo "$evaluate" | awk '{print toupper(substr($0,1,1))""substr($0,2)}' | sed -e 's/^/* /'` && evaluate=$'### :mag: About Evaluate\n\n'"${evaluate}"
[[ $tool = *[^[:space:]]* ]] && tool=`echo "$tool" | awk '{print toupper(substr($0,1,1))""substr($0,2)}' | sed -e 's/^/* /'` && tool=$'### :toolbox: Tools\n\n'"${tool}"
[[ $other = *[^[:space:]]* ]] && other=`echo "$other" | awk '{print toupper(substr($0,1,1))""substr($0,2)}' | sed -e 's/^/* /'` && other=$'### :sparkles: Others\n\n'"${other}"

RELEASE_BODY=$(cat <<- EOF
Release Date: $DATE_STRING

## :scroll: Changelog

$exploit

$evaluate

$tool

$other

## :key: Hash Table

|SHA256|EXECTUE FILE|
|---|---|
|$SHA256_TEXT_BODY|
EOF
)

TITLE="CDK $TAG_VERSION"

if [[ -z "${RELEASE_URL}" ]]; then
    echo -n "# $TITLE" $'\n\n' "$RELEASE_BODY" > /tmp/changelog.md 

else

    # update to github release page by github action
    RELEASE_BODY=`echo "$RELEASE_BODY" | jq -sR .`

    curl \
    -XPATCH \
    -H "${API_HEADER}" \
    -H "${AUTH_HEADER}" \
    -H "Content-Type: application/json" \
    -d "{\"name\": \"$TITLE\",\"body\": ${RELEASE_BODY}}" \
    "${RELEASE_URL}";

fi