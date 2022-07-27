#!/bin/sh

while getopts o: flag
do
    case "${flag}" in
        o) swagger=${OPTARG};;
    esac
done


routes='./temp-uploads/routes.yml'
output='./temp-uploads/output.yml'
temp='./temp-uploads/endpointData.yml'
methods='./temp-uploads/methods.yml'
pathAndParams='./temp-uploads/pathAndParams.yml'

cp $swagger $routes

echo 'endpoints:' > $output # begin output file

sed -i '/^  ".*":$/!d' $routes # store all paths

##################################

line1=$(head -n 1 $routes) # get first line
sed -i "1d" $routes # remove it from source

while ! [ -z $line1 ]; do # while $line1 not empty
    line2=$(head -n 1 $routes) # get first line
    sed -i "1d" $routes # remove it from source
    if [ -z $line2 ]; then # if methods.yml is empty...
        # echo line 1: $line1
        # echo "'til the end"

        tail -n +`grep -n -m 1 $line1 $swagger | cut -f1 -d:` $swagger > $temp
    else
        # echo line 1: $line1
        # echo line 2: $line2

        grep -A100000 $line1 $swagger | grep -B100000 $line2 > $temp # copy 1 path data
        sed -i '$d' $temp # remove the last line
    fi

    grep -w "get\|post\|delete\|parameters" $temp > $methods # store 2n indents (for counting)

    # echo lines in methods.yml:
    # cat $methods
    # echo ---

    # 2) existe-t-il la propriété "parameters" ?
    if grep -q 'parameters:' $methods; then
        hasParams='true'
        sed -i "1d" $methods # supprimer la premiere ligne de methods.yml ("parameters")
    else
        hasParams='false'
    fi

    cunt=$(< "./temp-uploads/methods.yml" wc -l) # count real amount of methods
    # echo "methods count:" $cunt

    # 5) pose un tag %pathandparams% avant les methods
    sed -i 's/get:$/%pathandparams%\n    get:/' $temp
    sed -i 's/post:$/%pathandparams%\n    post:/' $temp
    sed -i 's/delete:$/%pathandparams%\n    delete:/' $temp

    # on copie la route dans pathAndParams.yml
    head -1 $temp > $pathAndParams
    
    # if [ $hasParams == 'true' ]; then
    #     echo 'has parameters'
    # else
    #     echo 'no parameters'
    # fi

    # if [ \( $hasParams == 'true' \) -a \( $cunt -eq 1 \) ]; then
    #     echo '1 method, hasParams=true'
    # elif [ \( $hasParams == 'false' \) -a \( $cunt -eq 1 \) ]; then
    #     echo '1 method, hasParams=false'
    # elif [ \( $hasParams == 'true' \) -a \( $cunt -gt 1 \) ]; then
    #     echo 'multiple methods, hasParams=true'
    # elif [ \( $hasParams == 'false' \) -a \( $cunt -gt 1 \) ]; then
    #     echo 'multiple methods, hasParams=false'
    # fi
    # echo ---

    if [ $hasParams = 'false' ]; then
        # on fabrique un empty param.
        cp $temp tmp_file.yml
        awk 'NR==2{print "    parameters: []"}1' tmp_file.yml > $temp # add fake param to endpointData AVANT %pathandparams%

        if [ $cunt -eq 1 ]; then 
            # s'il n'y a qu'une méthode, on supprime l'unique %pathandparams% de endpointData (vu que path et (fake?)params sont déjà là)
            awk '/%pathandparams%/ && !f{f=1; next} 1' $temp > tmp_file.yml && mv tmp_file.yml $temp
        fi
    fi

    if [ $cunt -gt 1 ]; then
        # copier path+(fake?)parameters from endpointData.yml and save it to pathAndParams.yml
        awk "1;/%pathandparams%/{exit}" $temp > $pathAndParams

        echo  >> endpointData.yml # necessary to prevent disparition of last line
        # virer la derniere ligne (car c'est le nom de la méthode suivante)
        sed -i '$d' $pathAndParams
        # on rajoute le nom de propriété "path" devant la route de pathAndParams.yml, en en faisant un item d'array au passage, et en trimant le colon à la fin 
        sed -i '1 s/  [\S]*/  \- path: &/g' $pathAndParams
        sed -i '1 s/:$//' $pathAndParams



        # cat $pathAndParams # works


        # on supprime le tout premier %pathandparams% de endpointData (vu que path et (fake?)params sont déjà là)
        awk '/%pathandparams%/ && !f{f=1; next} 1' $temp > tmp_file.yml && mv tmp_file.yml $temp

        # TODO: improve
        # replace toutes les occurences de %pathandparams% par les path & param
        
        # in fact: put BEFORE %pathandparams% ...
        # sed -i 's/    %pathandparams%/tamere/' $temp # works

        # echo '############'
        # echo 'replace all %pathandparams% occurences:'
        sed "/%pathandparams%/{
        r $pathAndParams
        d
        }" $temp > tmp_file.yml && mv tmp_file.yml $temp
        
        
        # cat $temp
        # echo '############'
        
        

        # echo '####'
        # cat $temp
        # echo '####'
    fi
    # ...then delete every remaining %pathandparams%
    sed '/%pathandparams%/d' $temp > tmp_file.yml && mv tmp_file.yml $temp

    # on rajoute le nom de propriété "path" devant la premiere route de endpointData, en en faisant un item d'array au passage, et en trimant le colon à la fin
    sed -i '1 s/  [\S]*/  \- path: &/g' $temp
    sed -i '1 s/:$//' $temp

    # turn methods props into values
    sed -i 's/get:$/method:\n      type: GET/' $temp
    sed -i 's/post:$/method:\n      type: POST/' $temp
    sed -i 's/delete:$/method:\n      type: DELETE/' $temp

    cat $temp >> $output

    # echo '' # cute delimiter
    line1=$line2
done

sed -i -e '/      responses:$/{' -e 'n;s/.*/          code: &/' -e '}' $output # append "code:" before every https status code
sed -i 's/^\(          code:.*\):$/\1/' $output # remove trailing colon


echo '' >> $output

# clean up
rm $routes
rm $methods
rm $temp
rm $pathAndParams

echo 'parsing completed'
# echo '---'
# cat $output