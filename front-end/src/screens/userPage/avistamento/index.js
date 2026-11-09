import { View, Text, TouchableOpacity, TextInput, Image, Modal, ActivityIndicator } from "react-native";
import { SafeAreaView, } from "react-native-safe-area-context";
import { useNavigation } from "@react-navigation/native";
import Toast from "react-native-toast-message";
import { Feather } from "@expo/vector-icons"
import { useState } from 'react';
import styled from "./styled";

import * as Api from "../../../../services/API"

export default () => {

   
    const [loading, setLoading] = useState(false);

    

   

    return (
        <SafeAreaView style={styled.safeArea}>
                

            {loading && (
                <View style={styled.loadingOverlay}>
                    <ActivityIndicator size={'large'} color="#0000ff" />
                </View>
            )}
            <Modal
                transparent={true}
                animationType='none'
                visible={loading}
            />

        </SafeAreaView>
    );
};