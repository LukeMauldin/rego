(function() {

    'use strict';

    //Create application
    var regoApp = angular.module('regoApp', []);

    //Create controller
    regoApp.controller('mainCtrl', ['$scope', '$http', '$location', function($scope, $http, $location) {
        //Create emtpy matches
        $scope.matches = [];

        //Set default text
        $scope.regexpInput = 'r([a-z]+)go';
        $scope.stringInput = 'rego';
        $scope.findAllSubmatch = true;

        $scope.clearMatchResults = function() {
            //Clear all match output
            $scope.error = '';
            $scope.matches = [];
            $scope.matchResult = '';
        };

        $scope.clearAllFields = function() {
            $scope.regexpInput = '';
            $scope.stringInput = '';
            $scope.findAllSubmatch = true;
            $scope.clearMatchResults();
        };

        $scope.evaluateRegex = function() {
            //Retrieve updated regexp information
            var postData = {
                Regexp: $scope.regexpInput,
                Text: $scope.stringInput,
                FindAllSubmatch: $scope.findAllSubmatch === 'true' || $scope.findAllSubmatch === true,
            };
            var uri = $location.protocol() + "://" + $location.host() + ":" + $location.port() + "/test_regexp/";
            $http.post(uri,  postData)
                .success(function(data) {
                    //Clear results
                    $scope.clearMatchResults();

                    //Check for results
                    if (data.matches === null) {
                        return;
                    }

                    //Populate new results
                    var fullMatches = [];
                    for (var i = 0; i < data.matches.length; i++) {
                        var match = data.matches[i];

                        //Populate fullMatches list
                        fullMatches[i] = match[0];

                        for (var j = 1; j < match.length; j++) {
                            //Populate matches list
                            var groupName = data.groupsName.length >= j && data.groupsName[j - 1] !== '' ? data.groupsName[j - 1] : '-';
                            $scope.matches[j - 1] ={
                                    count: j - 1,
                                    groupName: groupName,
                                    matchText: match[j]
                            };
                        }
                    }
                    $scope.matchResult = fullMatches.join(" ");
                })
                .error(function(data) {
                    //Clear results
                    $scope.clearMatchResults();

                    //Populate error
                    $scope.error = data;
                });
        };

        //Invoke evaluateRegex to display initial data to user
        $scope.evaluateRegex();
    }]);

})();
