//(function() {

    'use strict';

    //Create application
    var regoApp = angular.module('regoApp', []);

    //Create controller
    regoApp.controller('mainCtrl', ['$scope', '$http', '$location', function($scope, $http, $location) {
        //Set default text
        $scope.regexpInput = 'r([a-z]+)go';
        $scope.stringInput = 'rego';
        $scope.findAllSubmatch = 'true';

        //Create emtpy matches
        $scope.matches = [];
//        $scope.matches = [
 //           {count: '0', 'groupName': '-', 'matchText': 'a match'}
   //     ];



        $scope.evaluateRegex = function() {
            //Clear all match output
            $scope.error = '';
            $scope.matches = [];
            $scope.matchResult = '';


            //Retrieve updated regexp information
            var postData = {
                Regexp: $scope.regexpInput,
                Text: $scope.stringInput,
                FindAllSubmatch: $scope.findAllSubmatch === 'true',
            };
            var uri = $location.protocol() + "://" + $location.host() + ":" + $location.port() + "/test_regexp/";
            $http.post(uri,  postData)
                .success(function(data) {
                    var fullMatches = [];
                    for (var i = 0; i < data.matches.length; i++) {
                        var match = data.matches[i];

                        //Populate fullMatches list
                        fullMatches[i] = match[0];

                        for (var j = 1; j < match.length; j++) {
                            //Populate matches list
                            var groupName = data.groupsName.length >= j && data.groupsName[j - 1] !== '' ? data.groupsName[j - 1] : '-';
                            $scope.matches[j] ={
                                    count: j,
                                    groupName: groupName,
                                    matchText: match[j]
                            };
                        }
                    }
                    $scope.matchResult = fullMatches.join(" ");
                })
                .error(function(data) {
                    $scope.error = data;
                });
        };

        //Invoke evaluateRegex to display initial data to user
        $scope.evaluateRegex();
    }]);

//})();
