import { View, Text, TouchableOpacity, TextInput, Image, ActivityIndicator, Modal, ImageBackground } from "react-native";
import { SafeAreaView, } from "react-native-safe-area-context";
import { useNavigation } from "@react-navigation/native";
import Toast from "react-native-toast-message";
import { Feather } from "@expo/vector-icons"
import { useState } from 'react';
import styled from "./styled";

export default () => {
    return (
        <SafeAreaView style={styled.safeArea}>
            <TouchableOpacity style={styled.cardContainer} onPress={''} activeOpacity={0.8}>
                <ImageBackground source={{ uri: "https://www.infoescola.com/wp-content/uploads/2018/03/sagui-de-tufos-brancos_790108426-1000x667.jpg" }} style={styled.cardBackground} imageStyle={styled.imageStyle} >
                    <View style={styled.overlay}>
                        <Text style={styled.cardText}>{"Meu Card"}</Text>
                    </View>
                </ImageBackground>
            </TouchableOpacity>
        </SafeAreaView>

)};